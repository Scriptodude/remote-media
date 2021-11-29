package mediahandler

import (
	log "github.com/sirupsen/logrus"

	"github.com/micmonay/keybd_event"
)

type MediaHandler interface {
	// PlayNext must play the next song, the first media player may be chosen
	PlayNext()

	// PlayPrevious must play the previous song, the first media player may be chosen
	PlayPrevious()

	// VolumeUp increases the volume of the media by a reasonable amount
	VolumeUp()

	// VolumeDown decreases the volume of the media by a reasonable amount
	VolumeDown()
}

// KeyboardMediaHandler uses keyboard events (/dev/uinput on linux for instance) to send events to the system
// This requires that the application have access to the keyboard controls programmatically
type KeyboardMediaHandler struct {
	kb  keybd_event.KeyBonding
	log *log.Entry
}
