#!/bin/bash

echo "STARTING SETUP..."

## add user to input ##
sudo usermod -a -G input user

## udev ##
sudo cp *.rules /etc/udev/rules.d

echo "REBOOTING SERVER..." && sudo shutdown -r now