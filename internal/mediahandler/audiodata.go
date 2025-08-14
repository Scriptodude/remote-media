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

	// TODO: Extract this part in a function called "findFirstRunningSink"
	for _, d := range data {
		log.Debugf("sink %+v", d)
		if d.state == "RUNNING" { //< TODO: Use a strong type instead of string comparison ?
			extractRegex, err := regexp.Compile("[0-9]{1,3}%") //< TODO: Extract the regex in a const

			if err != nil {
				log.Printf("unable to get volume, ignoring. %e", err)
				return -1
			}

			result := extractRegex.Find([]byte(d.volume))
			log.Debugf("regex result: %s", result)
			percent, err := strconv.Atoi(strings.Trim(string(result), "%"))

			if err != nil {
				log.Printf("unable to get volume, ignoring. %e", err)
				return -1
			}

			return percent
		}
	}

	return -1
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

	for lineNo := 0; lineNo < len(lines); lineNo++ {

		// Closure that takes a line relative to lineNo (lineNo+i),
		// removes the prop name and returns whatever data it contained, trimed of extra spaces
		cleanData := func(i int, propName string) string {
			if !strings.Contains(lines[lineNo+i], propName) {
				log.Warnf("Wrong line number for %s (%d+%d) - %s", propName, lineNo, i, lines[lineNo+i])
				return ""
			}

			dataNoName := strings.TrimLeft(strings.TrimSpace(lines[lineNo+i]), propName)
			return strings.TrimSpace(dataNoName)
		}

		// We start a new sink
		if strings.Contains(lines[lineNo], "Sink #") {
			state := cleanData(1, "State:")
			name := cleanData(2, "Name:")
			volume := cleanData(9, "Volume:")
			sink := sinkData{
				sinkName: name,
				state:    state,
				volume:   volume,
			}
			sinks = append(sinks, sink)
			log.Debugf("new sink %+v", sink)

			lineNo += 9
		}
	}

	if len(sinks) == 0 {
		return nil, errors.New("No sink found")
	}

	return sinks, nil
}
