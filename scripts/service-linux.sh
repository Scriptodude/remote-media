#!/bin/bash
# usage:
# chmod a+x ./service-linux.sh
# sudo ./service-linux.sh $USER

cd ..
dest="$(pwd)/bin"
mkdir -p "$dest"
/usr/local/go/bin/go build -o "$dest/remote-media" .
chmod +x "$dest/remote-media"
chown $1:$1 "$dest/remote-media"

cat > "/etc/systemd/system/remote-media.service" << EOM
[Unit]
Description=Remote media handler
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
User=$1
ExecStart="$dest/remote-media" -port=1337

[Install]
WantedBy=multi-user.target
EOM

systemctl enable "remote-media"
systemctl start "remote-media"