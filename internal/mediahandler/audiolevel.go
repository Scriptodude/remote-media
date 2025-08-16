package mediahandler

import (
	"github.com/scriptodude/remote-media/internal/pactl"
	"github.com/sirupsen/logrus"
)

const (
	// Matches volume in percent such as 100%, 49% or 0%
	extractVolumeRegex = "[0-9]{1,3}%"
)

type audioLevel struct {
	log        *logrus.Entry
	pulseAudio pactl.PulseAudioController
	pulseInfo  *pactl.PulseAudioInfo
}

// getVolumeLevel Returns the current audio level 0-100% of the first active audio sink
// using the pactl unix command
func (al *audioLevel) getVolumeLevel() int {
	al.init()

	level, err := al.pulseAudio.GetVolumeLevel(al.pulseInfo.DefaultSinkName)
	if err != nil {
		al.log.Errorln("Error getting volume level", err)
		return 0
	}

	return level
}

// setVolumeLevel sets the volume level of the default sink in percentage
func (al *audioLevel) setVolumeLevel(level int) {
	al.init()

	err := al.pulseAudio.SetVolumeLevel(al.pulseInfo.DefaultSinkName, level)
	if err != nil {
		al.log.Errorln("Error setting volume level", err)
	}
}

// Init initializes - if needed - the audio level struct with the required info for its good functioning.
func (al *audioLevel) init() {
	if al.pulseInfo == nil {
		al.log.Infoln("Fetching pulseInfo")
		info, err := al.pulseAudio.Info()

		if err != nil {
			al.log.Errorln("Error fetching pulse info", err)
			return
		}

		al.pulseInfo = &info
	}
}
