package mediahandler

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	// Matches volume in percent such as 100%, 49% or 0%
	extractVolumeRegex = "[0-9]{1,3}%"
)

type sinkData struct {
	sinkName string
	state    string
	volume   string
}

type audioLevel struct {
	previousActive *sinkData
	log            *logrus.Entry
}

// getAudioLevel Returns the current audio level 0-100% of the first active audio sink
// using the pactl unix command
func (al audioLevel) getAudioLevel() int {
	defaultLevel := al.log.Logger.Level
	al.log.Logger.Level = logrus.DebugLevel
	defer func() { al.log.Logger.Level = defaultLevel }()

	al.findRunningSink()

	if al.previousActive != nil {
		al.log.Debugln("Using existing sink ", al.previousActive.sinkName)
		volume, err := al.extractVolume()
		if err != nil {
			al.log.Printf("unable to get volume, ignoring. %e", err)
			return -1
		}

		return volume
	}

	al.log.Printf("unable to get volume, ignoring")
	return -1
}

func (al *audioLevel) setVolumeLevel(level int) {
	defaultLevel := al.log.Logger.Level
	al.log.Logger.Level = logrus.DebugLevel
	defer func() { al.log.Logger.Level = defaultLevel }()

	al.findRunningSink()

	if al.previousActive == nil {
		al.log.Warningln("No previous active sink")
		return
	}

	cmd := exec.Command("/usr/bin/pactl", "set-sink-volume", al.previousActive.sinkName, fmt.Sprintf("%d%%", level))
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	err := cmd.Run()
	if err != nil {
		al.log.Errorln("Cannot set volume level", err)
		return
	}
}

func (al *audioLevel) findRunningSink() {
	cmd := exec.Command("/usr/bin/pactl", "list", "sinks")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		al.log.Printf("unable to get volume, ignoring. %e", err)
		return
	}

	if err := cmd.Start(); err != nil {
		al.log.Printf("unable to get volume, ignoring. %e", err)
		return
	}

	data, err := parseData(al.log, stdout)
	if err != nil {
		al.log.Printf("unable to get volume, ignoring. %e", err)
		return
	}

	for _, d := range data {
		if d.state == "RUNNING" {
			al.previousActive = &d
			return
		}
	}

	al.previousActive = nil
}

// extractVolume finds the first sinkData with a state of RUNNING and returns its volume percentage
// In case of error, -1 is returned as the first return value and an error is also returned
func (al *audioLevel) extractVolume() (int, error) {
	if al.previousActive == nil {
		return -1, errors.New("No previous active sink")
	}

	extractRegex, err := regexp.Compile(extractVolumeRegex)
	if err != nil {
		return -1, err
	}

	result := extractRegex.FindString(al.previousActive.volume)
	percent, err := strconv.Atoi(strings.Trim(result, "%"))

	if err != nil {
		return -1, err
	}

	return percent, nil
}

// parseData parses the response coming from pactl into a sinkData struct
func parseData(log *logrus.Entry, reader io.Reader) ([]sinkData, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	stringData := string(data)
	lines := strings.Split(stringData, "\n")

	sinks := make([]sinkData, 0)
	lenOfLines := len(lines)
	for lineNo := 0; lineNo < lenOfLines; lineNo++ {

		// Closure that finds the first line relative to lineNo that contains the propName
		// returns empty string if none was found
		// removes the prop name and returns whatever data it contained, trimed of extra spaces
		findAndCleanData := func(propName string) string {
			for offset := 0; lineNo+offset < lenOfLines; offset++ {
				if !strings.Contains(lines[lineNo+offset], propName) {
					continue
				}

				dataNoName := strings.TrimLeft(strings.TrimSpace(lines[lineNo+offset]), propName)
				return strings.TrimSpace(dataNoName)
			}

			log.Warnf("%s was not found starting at line %d", propName, lineNo)
			return ""
		}

		// We have found a new sink - parse its data
		if strings.Contains(lines[lineNo], "Sink #") {
			state := findAndCleanData("State:")
			name := findAndCleanData("Name:")
			volume := findAndCleanData("Volume:")
			sink := sinkData{
				sinkName: name,
				state:    state,
				volume:   volume,
			}
			sinks = append(sinks, sink)
			log.Debugf("new sink %+v", sink)
		}
	}

	if len(sinks) == 0 {
		return nil, errors.New("No sink found")
	}

	return sinks, nil
}
