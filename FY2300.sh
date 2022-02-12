#! /bin/bash

/usr/bin/daemon -u pi:dialout -n frequency-generatord2 \
--env='HOME=/home/pi/git/frequency-server' \
--env='USBSPEED=9600' \
--env='USBPORT=usb-Silicon_Labs_CP2102_USB_to_UART_Bridge_Controller_0001-if00-port0' \
--stdout=/tmp/frequency-generator2.log \
/home/pi/git/frequency-server/frequency-server