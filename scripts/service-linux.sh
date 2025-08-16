#!/bin/bash
# usage:
# chmod a+x ./service-linux.sh
# ./service-linux.sh

APP_PATH=$HOME/.local/remote-media
SERVICE_PATH=$HOME/.config/systemd/user/remote-media.service

echo "Cleaning up previous stuff"
systemctl disable "remote-media"
systemctl stop "remote-media"
rm -f $SERVICE_PATH
rm -rf $APP_PATH

echo "Creating needed stuff"
mkdir -p $APP_PATH/web
cp -rR ../web/release $APP_PATH/web

echo "BUILDING"
cd ../cmd/remotemedia
/usr/local/go/bin/go build -o $APP_PATH/remote-media .

# TODO: No longer using dbus... this might no longer be needed
# if test -f "/etc/udev/rules.d/99-input.rules"; then
#     echo "No rules to write"
# else
#     echo 'KERNEL=="uinput", GROUP="uinput", MODE:="0660"' > "/etc/udev/rules.d/99-input.rules"
# fi

chmod 744 $APP_PATH/remote-media
chown $USER:root $APP_PATH/

echo "Creating service"
cat > $SERVICE_PATH << EOM
[Unit]
Description=Remote media handler
StartLimitIntervalSec=0
[Service]
Type=idle
Restart=always
RestartSec=15
WorkingDirectory=$APP_PATH/
ExecStart="$APP_PATH/remote-media" -port=1337

[Install]
WantedBy=default.target
EOM

echo "Enabled / starting service"
systemctl --user enable "remote-media"
systemctl --user restart "remote-media"
systemctl --user daemon-reload
