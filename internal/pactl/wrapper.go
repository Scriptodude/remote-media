package pactl

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/scriptodude/remote-media/internal/log"
	"github.com/sirupsen/logrus"
)

const (
	// Matches volume in percent such as 100%, 49% or 0%
	extractVolumeRegex = "[0-9]{1,3}%"
	pactlCmd           = "/usr/bin/pactl"
)

var DefaultCommandExecutor CommandExecutor = func(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	return cmd.Run()
}

var DefaultCommandExecutorWithStdout CommandExecutorWithStdout = func(command string, arg ...string) ([]byte, error) {
	cmd := exec.Command(command, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return data, nil
}

type pactlWrapper struct {
	executor           CommandExecutor
	executorWithStdout CommandExecutorWithStdout
	log                *logrus.Entry
}

func NewPulseAudioController(executor CommandExecutor, executorWithStdout CommandExecutorWithStdout) PulseAudioController {
	return &pactlWrapper{
		executor:           executor,
		executorWithStdout: executorWithStdout,
		log:                log.GetLoggerForHandler("PulseAudio"),
	}
}

// GetVolumeLevel implements PulseAudioController.
func (p *pactlWrapper) GetVolumeLevel(sink string) (int, error) {
	p.log.Infof("Fetching level for %s", sink)
	stdout, err := p.executorWithStdout(pactlCmd, "get-sink-volume", sink)
	if err != nil {
		return -1, err
	}

	extractRegex, err := regexp.Compile(extractVolumeRegex)
	if err != nil {
		return -1, err
	}

	result := extractRegex.Find(stdout)
	percent, err := strconv.Atoi(strings.Trim(string(result), "%"))
	if err != nil {
		return -1, err
	}

	return percent, nil
}

// Info implements PulseAudioController.
func (p *pactlWrapper) Info() (PulseAudioInfo, error) {
	stdout, err := p.executorWithStdout(pactlCmd, "--format=json", "info")
	if err != nil {
		return PulseAudioInfo{}, err
	}

	var info PulseAudioInfo
	err = json.Unmarshal(stdout, &info)

	return info, err
}

// SetVolumeLevel implements PulseAudioController.
func (p *pactlWrapper) SetVolumeLevel(sink string, level int) error {
	return p.executor(pactlCmd, "set-sink-volume", sink, fmt.Sprintf("%d%%", level))
}
