package remotemedia_test

import (
	"bytes"
	"testing"

	"github.com/scriptodude/remote-media/internal/remotemedia"
)

type mock_media struct {
	volumeUp     func() int
	volumeDown   func() int
	setVolume    func(int) int
	getVolume    func() int
	playNext     func()
	playPrevious func()
}

func TestGetCurrentState_TakesVolumeFromMediaHandler(t *testing.T) {
	m := mock_media{
		getVolume: func() int {
			return 15
		},
	}

	broker := remotemedia.NewMessageBroker(m)

	result := broker.GetCurrentState()

	if !bytes.Equal(result, []byte("15")) {
		t.Errorf("Wrong value, expect 15, got %s", result)
	}
}

func TestHandleMessage_ShouldReturnNil_WhenOpCodeIsntText(t *testing.T) {
	m := mock_media{}

	broker := remotemedia.NewMessageBroker(m)

	result := broker.HandleMessage(0x25, []byte{})

	if result != nil {
		t.Errorf("Wrong value, expect nil, got %s", result)
	}
}

func TestHandleMessage_ShouldReturnNil_WhenSetVolumeFormatIsInvalid(t *testing.T) {
	m := mock_media{}

	broker := remotemedia.NewMessageBroker(m)

	result := broker.HandleMessage(0x1, []byte("set_volume==5"))

	if result != nil {
		t.Errorf("Wrong value, expected nil, got %s", result)
	}
}

func TestHandleMessage_ShouldReturnNil_WhenSetVolumeNumberFormatIsInvalid(t *testing.T) {
	m := mock_media{}

	broker := remotemedia.NewMessageBroker(m)

	result := broker.HandleMessage(0x1, []byte("set_volume=a"))

	if result != nil {
		t.Errorf("Wrong value, expected nil, got %s", result)
	}
}

func TestHandleMessage_ShouldReturnSetVolumeData_WhenMessageIsSetVolume(t *testing.T) {
	m := mock_media{
		setVolume: func(v int) int { return v },
	}

	broker := remotemedia.NewMessageBroker(m)

	result := broker.HandleMessage(0x1, []byte("set_volume=5"))

	if !bytes.Equal(result, []byte("5")) {
		t.Errorf("Wrong value, expected 5, got %s", result)
	}
}

func TestHandleMessage_ShouldReturnVolumeUpData_WhenMessageIsVolumeUp(t *testing.T) {
	m := mock_media{
		volumeUp: func() int { return 5 },
	}

	broker := remotemedia.NewMessageBroker(m)

	result := broker.HandleMessage(0x1, []byte("volume_up"))

	if !bytes.Equal(result, []byte("5")) {
		t.Errorf("Wrong value, expected 5, got %s", result)
	}
}

func TestHandleMessage_ShouldReturnVolumeUpData_WhenMessageIsVolumeDown(t *testing.T) {
	m := mock_media{
		volumeDown: func() int { return 5 },
	}

	broker := remotemedia.NewMessageBroker(m)

	result := broker.HandleMessage(0x1, []byte("volume_down"))

	if !bytes.Equal(result, []byte("5")) {
		t.Errorf("Wrong value, expected 5, got %s", result)
	}
}

func TestHandleMessage_ShouldCallPlayNext_WhenMessageIsPlayNext(t *testing.T) {
	count := 0
	m := mock_media{
		playNext: func() { count = 1 },
	}

	broker := remotemedia.NewMessageBroker(m)

	broker.HandleMessage(0x1, []byte("play_next"))

	if count != 1 {
		t.Errorf("expected play next to be called")
	}
}

func TestHandleMessage_ShouldCallPlayPrevious_WhenMessageIsPlayPrevious(t *testing.T) {
	count := 0
	m := mock_media{
		playPrevious: func() { count = 1 },
	}

	broker := remotemedia.NewMessageBroker(m)

	broker.HandleMessage(0x1, []byte("play_previous"))

	if count != 1 {
		t.Errorf("expected play next to be called")
	}
}

// implementation of mock
func (m mock_media) PlayNext() {
	m.playNext()
}

// PlayPrevious must play the previous song, the first media player may be chosen
func (m mock_media) PlayPrevious() {
	m.playPrevious()
}

// VolumeUp increases the volume of the media by a reasonable amount and returns the current volume level (0-100)
func (m mock_media) VolumeUp() int {
	return m.volumeUp()
}

// VolumeDown decreases the volume of the media by a reasonable amount and returns the current volume level (0-100)
func (m mock_media) VolumeDown() int {
	return m.volumeDown()
}

// SetVolume sets the volume to the given percentage
func (m mock_media) SetVolume(v int) int {
	return m.setVolume(v)
}

// GetVolume gets the current volume percentage
func (m mock_media) GetVolume() int {
	return m.getVolume()
}
