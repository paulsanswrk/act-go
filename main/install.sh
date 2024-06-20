#!/bin/bash

cat > /etc/systemd/system/act.service <<EOF
[Unit]
Description=ACT Service
After=network.target
[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/act
ExecStart=/home/ubuntu/act/act.bin
Restart=always
[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable act
sudo systemctl start act