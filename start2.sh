#! /bin/bash

/usr/bin/daemon -u pi:dialout -n frequency-serverd2 \
--env='HOME=/home/pi/git/frequency-server' \
--env='WEBPORT=8082' \
--env='USBSPEED=9600' \
--env='USBPORT=usb-Silicon_Labs_CP2102_USB_to_UART_Bridge_Controller_0001-if00-port0' \
--stdout=/tmp/frequency-server2.log \
/home/pi/git/frequency-server/frequency-server
