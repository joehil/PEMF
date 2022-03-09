#! /bin/bash

/usr/bin/nice -n -5 /usr/bin/daemon -u pi:dialout -n PEMF-generatord2 \
--env='HOME=/home/pi/git/PEMF' \
--env='USBSPEED=9600' \
--env='USBPORT=usb-Silicon_Labs_CP2102_USB_to_UART_Bridge_Controller_0001-if00-port0' \
--env='PIPE=/tmp/FY2300.pipe' \
--stdout=/tmp/PEMF-generator2.log \
/home/pi/git/PEMF/PEMF generator
