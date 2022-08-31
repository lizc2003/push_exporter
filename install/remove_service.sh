#!/bin/bash

NAME=push_exporter.service

sudo cp ./$NAME /lib/systemd/system/
sudo systemctl stop $NAME
sudo systemctl disable $NAME
sudo rm -f /lib/systemd/system/$NAME
sudo systemctl daemon-reload
sudo systemctl reset-failed
