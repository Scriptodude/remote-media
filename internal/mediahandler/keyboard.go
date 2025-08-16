package mediahandler

import (
	"github.com/micmonay/keybd_event"
	"github.com/scriptodude/remote-media/internal/log"
	"github.com/scriptodude/remote-media/internal/pactl"

	"github.com/sirupsen/logrus"
)

// keyboardMediaHandler uses keyboard events (/dev/uinput on linux for instance) to send events to the system
// This requires that the application have access to the keyboard controls programmatically
type keyboardMediaHandler struct {
	kb  keybd_event.KeyBonding
	log *logrus.Entry
	al  audioLevel
}

func NewKeyboardMediaHandler() MediaHandler {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	log := log.GetLoggerForHandler("Keyboard")
	return &keyboardMediaHandler{
		kb:  kb,
		log: log,
		al: audioLevel{
			log:        log,
			pulseAudio: pactl.NewPulseAudioController(pactl.DefaultCommandExecutor, pactl.DefaultCommandExecutorWithStdout),
		},
	}
}

// PlayNext implements the MediaHandler interface and plays the next song
func (k *keyboardMediaHandler) PlayNext() {
	k.log.Info("Playing next song")
	tapKey(keybd_event.VK_NEXTSONG, k.kb)
}

// PlayPrevious implements the MediaHandler interface and plays the previous song
func (k *keyboardMediaHandler) PlayPrevious() {
	k.log.Info("Playing previous song")
	tapKey(keybd_event.VK_PREVIOUSSONG, k.kb)
}

// VolumeUp implements the MediaHandler interface and increases the volume
func (k *keyboardMediaHandler) VolumeUp() int {
	k.log.Info("Increasing Volume")
	tapKey(keybd_event.VK_VOLUMEUP, k.kb)
	return k.al.getVolumeLevel()
}

// VolumeUp implements the MediaHandler interface and increases the volume
func (k *keyboardMediaHandler) VolumeDown() int {
	k.log.Info("Decreasing volume")
	tapKey(keybd_event.VK_VOLUMEDOWN, k.kb)
	return k.al.getVolumeLevel()
}

// tapKey taps a key on the keyboard
func tapKey(key int, kb keybd_event.KeyBonding) {
	kb.SetKeys(key)
	kb.Press()
	kb.Release()
}

// GetVolume implements MediaHandler.
func (k *keyboardMediaHandler) GetVolume() int {
	return k.al.getVolumeLevel()
}

// SetVolume implements MediaHandler.
func (k *keyboardMediaHandler) SetVolume(level int) int {
	k.al.setVolumeLevel(level)
	return k.al.getVolumeLevel()
}
