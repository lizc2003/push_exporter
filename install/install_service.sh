#!/bin/bash

NAME=push_exporter.service

sudo cp ./$NAME /lib/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable $NAME
