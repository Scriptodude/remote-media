#!/bin/bash
# usage:
# chmod a+x ./service-linux.sh
# sudo ./service-linux.sh $USER

echo "Cleaning up previous stuff"
systemctl disable "remote-media"
systemctl stop "remote-media"
rm -f "/etc/systemd/system/remote-media.service"
rm -rf /usr/local/remote-media 

echo "Creating needed stuff"
mkdir -p /usr/local/remote-media
cp -rR ../static /usr/local/remote-media

echo "BUILDING"
cd ..
/usr/local/go/bin/go build -o /usr/local/remote-media/remote-media .
chmod 744 /usr/local/remote-media/remote-media
chmod +0666 /dev/uinput
chown $1:root /usr/local/remote-media/remote-media

echo "Creating service"
cat > "/etc/systemd/system/remote-media.service" << EOM
[Unit]
Description=Remote media handler
StartLimitIntervalSec=0
[Service]
Type=idle
Restart=always
RestartSec=15
User=$1
WorkingDirectory=/usr/local/remote-media/
ExecStart="/usr/local/remote-media/remote-media" -port=1337

[Install]
WantedBy=default.target
EOM

echo "Enabled / starting service"
systemctl enable "remote-media"
systemctl restart "remote-media"
systemctl daemon-reload