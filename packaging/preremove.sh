#!/bin/bash
set -e

if systemctl is-active --quiet go-template-httpserver; then
    systemctl stop go-template-httpserver
fi

if systemctl is-enabled --quiet go-template-httpserver; then
    systemctl disable go-template-httpserver
fi
