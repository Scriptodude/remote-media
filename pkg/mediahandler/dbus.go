package mediahandler

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
	"github.com/scriptodude/remote-media/internal/log"
)

const (
	playNextMethod     = "org.mpris.MediaPlayer2.Player.Next"
	playPreviousMethod = "org.mpris.MediaPlayer2.Player.Previous"
	volumeProperty     = "org.mpris.MediaPlayer2.Player.Volume"
	propertyInterface  = "org.freedesktop.DBus.Properties"
	mprisPath          = "/org/mpris/MediaPlayer2"
)

func NewDbusMediaHandler() *DbusMediaHandler {

	return &DbusMediaHandler{
		media:      nil,
		properties: nil,
		log:        log.GetLoggerForHandler("Dbus"),
	}
}

func (h *DbusMediaHandler) PlayNext() {
	h.callMethod(playNextMethod)
}

func (h *DbusMediaHandler) PlayPrevious() {
	h.callMethod(playPreviousMethod)
}

func (h *DbusMediaHandler) VolumeUp() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	prop, err := getVolumeProperties(conn)
	if err != nil {
		panic(err)
	}
	prop.Set("org.mpris.MediaPlayer2.Player", "Volume", dbus.MakeVariant(float64(100)))
}

func (h *DbusMediaHandler) VolumeDown() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	prop, err := getVolumeProperties(conn)
	if err != nil {
		panic(err)
	}
	// org.gnome.SettingsDaemon.Sound
	prop.Set("org.mpris.MediaPlayer2.Player", "Volume", dbus.MakeVariant(float64(-1)))
}

func (h *DbusMediaHandler) callMethod(method string) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var names []string
	conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)

	var selectedName string
	for _, name := range names {
		h.log.Debugf("name : %s, %b", name, strings.Contains(name, "mpris"))
		if strings.Contains(name, "mpris") {
			selectedName = name
		}
	}
	h.log.Infof("Selected name is %s", selectedName)

	conn.Object(selectedName, mprisPath).Call(playPreviousMethod, dbus.FlagNoReplyExpected, 0)
}

func getVolumeProperties(conn *dbus.Conn) (*prop.Properties, error) {
	propsSpec := map[string]map[string]*prop.Prop{
		"org.mpris.MediaPlayer2.Player": {
			"Volume": {
				Value:    float64(0),
				Writable: true,
				Emit:     prop.EmitTrue,
				Callback: func(c *prop.Change) *dbus.Error {
					fmt.Println(c.Iface, c.Name, "changed to", c.Value)
					return nil
				},
			},
		},
	}

	return prop.Export(conn, dbus.ObjectPath(mprisPath), propsSpec)
}
