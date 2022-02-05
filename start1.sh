#! /bin/bash

/usr/bin/daemon -u pi:dialout -n frequency-serverd1 \
--env='HOME=/home/pi/git/frequency-server' \
--env='WEBPORT=8080' \
--env='USBSPEED=115200' \
--env='USBPORT=usb-1a86_USB_Serial-if00-port0' \
--env='GENFACTOR=1' \
--stdout=/tmp/frequency-server1.log \
/home/pi/git/frequency-server/frequency-server
