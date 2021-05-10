#!/bin/bash
# usage:
# chmod a+x ./service-linux.sh
# sudo ./service-linux.sh $USER

systemctl disable "remote-media"

echo "Cleaning up previous stuff"
rm -f "/etc/systemd/system/remote-media.service"
rm -f /usr/local/remote-media 

echo "BUILDING"
cd ..
/usr/local/go/bin/go build -o /usr/local/remote-media .
chmod 744 scripts/remote-media

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
ExecStart="/usr/local/remote-media" -port=1337

[Install]
WantedBy=default.target
EOM

echo "Enabled / starting service"
systemctl enable "remote-media"
systemctl start "remote-media"
systemctl daemon-reload