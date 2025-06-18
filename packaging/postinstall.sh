#!/bin/bash
set -e

if ! id "go-template" &>/dev/null; then
    useradd --system --no-create-home --shell /bin/false go-template
fi

mkdir -p /var/log/go-template
chown go-template:go-template /var/log/go-template
chmod 755 /var/log/go-template

systemctl daemon-reload
systemctl enable go-template-httpserver.service

echo "go-template-httpserver installed successfully!"
echo "To start: sudo systemctl start go-template-httpserver"
