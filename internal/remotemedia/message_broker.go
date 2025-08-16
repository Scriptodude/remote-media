package remotemedia

import (
	"strconv"
	"strings"

	"github.com/scriptodude/remote-media/internal/log"
	"github.com/scriptodude/remote-media/internal/mediahandler"
	"github.com/sirupsen/logrus"
)

type MessageBroker interface {
	HandleMessage(opcode int, message []byte) []byte
	GetCurrentState() []byte
}

type messageBrokerImpl struct {
	log          *logrus.Entry
	mediahandler mediahandler.MediaHandler
}

func NewMessageBroker() MessageBroker {

	// TODO: Allow the caller to chose the mediahandler impl ?
	return &messageBrokerImpl{
		log:          log.GetLoggerForHandler("MessageBroker"),
		mediahandler: mediahandler.NewKeyboardMediaHandler(),
	}
}

func (mb *messageBrokerImpl) GetCurrentState() []byte {
	return []byte(strconv.Itoa(mb.mediahandler.GetVolume()))
}

// HandleMessage implements the MessageBroker interface
func (mb *messageBrokerImpl) HandleMessage(opcode int, message []byte) []byte {
	mb.log.Infof("Received message %d, %s", opcode, message)

	// Opcode for text
	if opcode != 0x1 {
		mb.log.Warnln("Not supporting message type, continuing")
		return nil
	}

	trimmed := strings.TrimSpace(string(message))

	if strings.Contains(trimmed, "set_volume") {
		split := strings.Split(trimmed, "=")

		if len(split) != 2 {
			mb.log.Errorf("Invalid format for set_volume %s", message)
			return nil
		}

		volumeLevel, err := strconv.Atoi(split[1])
		if err != nil {
			mb.log.Errorf("Invalid format for set_volume %s", message)
			return nil
		}

		return []byte(strconv.Itoa(mb.mediahandler.SetVolume(volumeLevel)))
	}

	switch trimmed {
	case "volume_up":
		return []byte(strconv.Itoa(mb.mediahandler.VolumeUp()))

	case "volume_down":
		return []byte(strconv.Itoa(mb.mediahandler.VolumeDown()))

	case "play_next":
		mb.mediahandler.PlayNext()

	case "play_previous":
		mb.mediahandler.PlayPrevious()

	default:
		mb.log.Errorf("Unkown command: %s", message)
	}

	return nil
}
