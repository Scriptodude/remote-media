package mediahandler

import (
	"errors"
	"testing"

	"github.com/scriptodude/remote-media/internal/pactl"
	"github.com/sirupsen/logrus"
)

type mock_ctrl struct {
	getVolumeLevel func() (int, error)
	setVolumeLevel func(int) error
	info           func() (pactl.PulseAudioInfo, error)
}

func (c mock_ctrl) GetVolumeLevel(string) (int, error) {
	if c.getVolumeLevel == nil {
		panic("undefined getVolumeLevel")
	}
	return c.getVolumeLevel()
}

func (c mock_ctrl) SetVolumeLevel(sink string, v int) error {
	if c.setVolumeLevel == nil {
		panic("undefined setVolumeLevel")
	}
	return c.setVolumeLevel(v)
}

func (c mock_ctrl) Info() (pactl.PulseAudioInfo, error) {
	if c.info == nil {
		panic("undefined info")
	}
	return c.info()
}

func TestInit_ShouldKeepPulseInfoToNil_WhenErrorFromInfo(t *testing.T) {
	m := mock_ctrl{
		info: func() (pactl.PulseAudioInfo, error) { return pactl.PulseAudioInfo{}, errors.New("failure") },
	}

	sut := audioLevel{
		log:        logrus.NewEntry(logrus.New()),
		pulseAudio: m,
	}

	sut.init()

	if sut.pulseInfo != nil {
		t.Errorf("Expected nil pulseinfo %+v", sut.pulseInfo)
	}
}

func TestInit_ShouldSetPulseInfo_WhenNoError(t *testing.T) {
	m := mock_ctrl{
		info: func() (pactl.PulseAudioInfo, error) { return pactl.PulseAudioInfo{}, nil },
	}

	sut := audioLevel{
		log:        logrus.NewEntry(logrus.New()),
		pulseAudio: m,
	}

	sut.init()

	if sut.pulseInfo == nil {
		t.Error("Expected pulseinfo got nil")
	}
}

func TestInit_ShouldDoNothing_WhenPulseInfo_IsDefined(t *testing.T) {
	m := mock_ctrl{
		info: func() (pactl.PulseAudioInfo, error) { panic("Should not be called") },
	}

	sut := audioLevel{
		log:        logrus.NewEntry(logrus.New()),
		pulseAudio: m,
		pulseInfo:  &pactl.PulseAudioInfo{},
	}

	sut.init()

	if sut.pulseInfo == nil {
		t.Error("Expected pulseinfo got nil")
	}
}

func TestGetVolumeLevel_ShouldReturn0_UponError(t *testing.T) {
	m := mock_ctrl{
		getVolumeLevel: func() (int, error) { return -1, errors.New("a") },
	}

	sut := audioLevel{
		log:        logrus.NewEntry(logrus.New()),
		pulseAudio: m,
		pulseInfo:  &pactl.PulseAudioInfo{},
	}

	result := sut.getVolumeLevel()

	if result != 0 {
		t.Errorf("Did not return 0, got %d", result)
	}
}

func TestGetVolumeLevel_ShouldReturnVolumeLevel_OfPulseAudio(t *testing.T) {
	m := mock_ctrl{
		getVolumeLevel: func() (int, error) { return 55, nil },
	}

	sut := audioLevel{
		log:        logrus.NewEntry(logrus.New()),
		pulseAudio: m,
		pulseInfo:  &pactl.PulseAudioInfo{},
	}

	result := sut.getVolumeLevel()

	if result != 55 {
		t.Errorf("Did not return 55, got %d", result)
	}
}

func TestSetVolumeLevel_ShouldCallPulseAudio_WithLevel(t *testing.T) {
	m := mock_ctrl{
		setVolumeLevel: func(level int) error {
			if level != 55 {
				t.Errorf("expected to be called with 55, got %d", level)
			}
			return nil
		},
	}

	sut := audioLevel{
		log:        logrus.NewEntry(logrus.New()),
		pulseAudio: m,
		pulseInfo:  &pactl.PulseAudioInfo{},
	}

	sut.setVolumeLevel(55)
}
