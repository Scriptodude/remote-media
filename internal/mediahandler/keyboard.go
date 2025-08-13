package mediahandler

import (
	"github.com/micmonay/keybd_event"
	"github.com/scriptodude/remote-media/internal/log"

	"github.com/sirupsen/logrus"
)

// keyboardMediaHandler uses keyboard events (/dev/uinput on linux for instance) to send events to the system
// This requires that the application have access to the keyboard controls programmatically
type keyboardMediaHandler struct {
	kb  keybd_event.KeyBonding
	log *logrus.Entry
}

func NewKeyboardMediaHandler() MediaHandler {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	return &keyboardMediaHandler{
		kb:  kb,
		log: log.GetLoggerForHandler("Keyboard"),
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
	tapKey(keybd_event.VK_NEXTSONG, k.kb)
}

// VolumeUp implements the MediaHandler interface and increases the volume
func (k *keyboardMediaHandler) VolumeUp() int {
	k.log.Info("Increasing Volume")
	tapKey(keybd_event.VK_VOLUMEUP, k.kb)
	return getAudioLevel(k.log)
}

// VolumeUp implements the MediaHandler interface and increases the volume
func (k *keyboardMediaHandler) VolumeDown() int {
	k.log.Info("Decreasing volume")
	tapKey(keybd_event.VK_VOLUMEDOWN, k.kb)
	return getAudioLevel(k.log)
}

// tapKey taps a key on the keyboard
func tapKey(key int, kb keybd_event.KeyBonding) {
	kb.SetKeys(key)
	kb.Press()
	kb.Release()
}
