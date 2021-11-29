package mediahandler

import (
	"github.com/micmonay/keybd_event"
	"github.com/scriptodude/remote-media/internal/log"
)

func NewKeyboardMediaHandler() *KeyboardMediaHandler {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	return &KeyboardMediaHandler{
		kb:  kb,
		log: log.GetLoggerForHandler("Keyboard"),
	}
}

func (k *KeyboardMediaHandler) PlayNext() {
	k.log.Info("Playing next song")
	change(keybd_event.VK_NEXTSONG, k.kb)
}

func (k *KeyboardMediaHandler) PlayPrevious() {
	k.log.Info("Playing previous song")
	change(keybd_event.VK_NEXTSONG, k.kb)
}

func (k *KeyboardMediaHandler) VolumeUp() {
	k.log.Info("Increasing Volume")
	change(keybd_event.VK_VOLUMEUP, k.kb)
}

func (k *KeyboardMediaHandler) VolumeDown() {
	k.log.Info("Decreasing volume")
	change(keybd_event.VK_VOLUMEDOWN, k.kb)
}

func change(key int, kb keybd_event.KeyBonding) {
	kb.SetKeys(key)
	kb.Press()
	kb.Release()
}
