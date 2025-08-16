package pactl

type PulseAudioController interface {
	// Info returns a PulseAudioInfo structure with the info extracted
	Info() (PulseAudioInfo, error)

	// GetVolumeLevel gets the volume level percentage of the given sink
	GetVolumeLevel(sink string) (int, error)

	// SetVolumeLevel sets the volume level percentage of the given sink
	SetVolumeLevel(sink string, level int) error
}

type PulseAudioInfo struct {
	DefaultSampleSpecification string `json:"default_sample_specification"`
	DefaultChannelMap          string `json:"default_channel_map"`
	DefaultSinkName            string `json:"default_sink_name"`
}

type CommandExecutor func(string, ...string) error

type CommandExecutorWithStdout func(string, ...string) ([]byte, error)
