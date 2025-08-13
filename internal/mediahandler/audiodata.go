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

	for _, d := range data {
		log.Debugf("sink %+v", d)
		if d.state == "RUNNING" {
			extractRegex, err := regexp.Compile("[0-9]{1,3}%")

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

func parseData(log *logrus.Entry, reader io.Reader) ([]sinkData, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	stringData := string(data)
	numberOfSink := strings.Count(stringData, "Sink #")
	if numberOfSink == 0 {
		return nil, errors.New("No sink found")
	}

	sinks := make([]sinkData, numberOfSink)
	currentIdx := 0
	for i := 0; i < numberOfSink; i++ {
		currentIdx = strings.Index(stringData, "Sink #")
		newLine := strings.Index(stringData[currentIdx:], "\n")
		name := stringData[currentIdx : currentIdx+newLine]

		currentIdx += newLine
		stateLocation := strings.Index(stringData[currentIdx:], "State:")
		newLine = strings.Index(stringData[currentIdx+stateLocation:], "\n")
		state := stringData[currentIdx+stateLocation : currentIdx+stateLocation+newLine]

		currentIdx += newLine
		volumeLocation := strings.Index(stringData[currentIdx:], "Volume:")
		newLine = strings.Index(stringData[currentIdx+volumeLocation:], "\n")
		log.Debugf("%d %d", volumeLocation, newLine)
		volume := stringData[currentIdx+volumeLocation : currentIdx+volumeLocation+newLine]

		sinks[i] = sinkData{
			sinkName: strings.TrimSpace(name),
			state:    strings.TrimSpace(strings.TrimLeft(state, "State:")),
			volume:   strings.TrimSpace(strings.TrimLeft(volume, "Volume:")),
		}

		log.Debugf("new sink %+v", sinks[i])

		currentIdx = newLine
		stringData = stringData[newLine:]
	}

	return sinks, nil
}
