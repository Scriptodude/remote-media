package mediahandler

type MediaHandler interface {
	// PlayNext must play the next song, the first media player may be chosen
	PlayNext()

	// PlayPrevious must play the previous song, the first media player may be chosen
	PlayPrevious()

	// VolumeUp increases the volume of the media by a reasonable amount and returns the current volume level (0-100)
	VolumeUp() int

	// VolumeDown decreases the volume of the media by a reasonable amount and returns the current volume level (0-100)
	VolumeDown() int
}
