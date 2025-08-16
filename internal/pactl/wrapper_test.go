package pactl

import (
	"errors"
	"testing"
)

func getExecutor(err error) CommandExecutor {
	return func(string, ...string) error {
		return err
	}
}

func getExecutorWithStdout(data string, err error) CommandExecutorWithStdout {
	return func(string, ...string) ([]byte, error) {
		if err != nil {
			return nil, err
		}

		return []byte(data), err
	}
}

func TestSetVolumeLevel_ShouldReturnErrorOfExecutor_WhenError(t *testing.T) {
	sut := NewPulseAudioController(getExecutor(errors.New("failure")), getExecutorWithStdout("", nil))
	err := sut.SetVolumeLevel("name", 1)

	if err == nil {
		t.Error("Expected error, found nil")
	}
}

func TestSetVolumeLevel_ShouldReturnNil_WhenNoError(t *testing.T) {
	sut := NewPulseAudioController(getExecutor(nil), getExecutorWithStdout("", nil))
	err := sut.SetVolumeLevel("name", 1)

	if err != nil {
		t.Errorf("Expected nil; found %e", err)
	}
}

func TestInfo_ShouldReturnError_WhenExecutorWithStdout_HasError(t *testing.T) {
	sut := NewPulseAudioController(getExecutor(nil), getExecutorWithStdout("", errors.New("failure")))
	_, err := sut.Info()

	if err == nil {
		t.Error("Expected error, found nil")
	}
}

func TestInfo_ShouldReturnError_WhenJsonCannotBeParsed(t *testing.T) {
	sut := NewPulseAudioController(getExecutor(nil), getExecutorWithStdout("", nil))
	_, err := sut.Info()

	if err == nil {
		t.Error("Expected error, found nil")
	}
}

func TestInfo_ShouldReturnParsedPulseAudioInfo_WhenNoError(t *testing.T) {
	sut := NewPulseAudioController(
		getExecutor(nil),
		getExecutorWithStdout(`{"default_sink_name":"a","default_sample_specification":"b","default_channel_map":"c"}`, nil),
	)
	result, err := sut.Info()

	if err != nil {
		t.Error("Expected nil, found ", err)
	}

	if result.DefaultSinkName != "a" || result.DefaultSampleSpecification != "b" || result.DefaultChannelMap != "c" {
		t.Errorf("result not as expected %v", result)
	}
}

func TestGetVolumeLevel_ShouldReturnErrorOfExecutor_WhenError(t *testing.T) {
	sut := NewPulseAudioController(getExecutor(nil), getExecutorWithStdout("", errors.New("failure")))
	_, err := sut.GetVolumeLevel("name")

	if err == nil {
		t.Error("Expected error, found nil")
	}
}

func TestGetVolumeLevel_ShouldReturnErrorOfAtoi_WhenNoNumberWasFound(t *testing.T) {
	sut := NewPulseAudioController(getExecutor(nil), getExecutorWithStdout("asdf%", nil))
	_, err := sut.GetVolumeLevel("name")

	if err == nil {
		t.Error("Expected error, found nil")
	}
}

func TestGetVolumeLevel_ShouldReturnVolumeLevel(t *testing.T) {
	sut := NewPulseAudioController(getExecutor(nil), getExecutorWithStdout("15%", nil))
	level, err := sut.GetVolumeLevel("name")

	if err != nil {
		t.Error("Expected nil, found ", err)
	}

	if level != 15 {
		t.Errorf("Expected 15, found %d", level)
	}
}
