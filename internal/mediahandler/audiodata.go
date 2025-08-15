package mediahandler

import (
	"errors"
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

// getAudioLevel Returns the current audio level 0-100% of the first active audio sink
// using the pactl unix command
func getAudioLevel(log *logrus.Entry) int {
	defaultLevel := log.Logger.Level
	log.Logger.Level = logrus.DebugLevel
	defer func() { log.Logger.Level = defaultLevel }()

	cmd := exec.Command("pactl", "list", "sinks")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("unable to get volume, ignoring. %e", err)
		return -1
	}

	if err := cmd.Start(); err != nil {
		log.Printf("unable to get volume, ignoring. %e", err)
		return -1
	}

	data, err := parseData(log, stdout)
	if err != nil {
		log.Printf("unable to get volume, ignoring. %e", err)
		return -1
	}

	volume, err := findFirstRunningSinkVolume(data)
	if err != nil {
		log.Printf("unable to get volume, ignoring. %e", err)
		return -1
	}

	return volume
}

type sinkData struct {
	sinkName string
	state    string
	volume   string
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

// findFirstRunningSinkVolume finds the first sinkData with a state of RUNNING and returns its volume percentage
// In case of error, -1 is returned as the first return value and an error is also returned
func findFirstRunningSinkVolume(sinks []sinkData) (int, error) {
	extractRegex, err := regexp.Compile(extractVolumeRegex)
	if err != nil {
		return -1, err
	}

	for _, d := range sinks {
		if d.state == "RUNNING" {
			result := extractRegex.FindString(d.volume)
			percent, err := strconv.Atoi(strings.Trim(result, "%"))

			if err != nil {
				return -1, err
			}

			return percent, nil
		}
	}

	return -1, errors.New("No running sink found")
}
