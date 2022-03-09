#! /bin/bash

/usr/bin/nice -n -5 /usr/bin/daemon -u pi:dialout -n PEMF-generatord1 \
--env='HOME=/home/pi/git/PEMF' \
--env='USBSPEED=115200' \
--env='USBPORT=usb-1a86_USB_Serial-if00-port0' \
--env='PIPE=/tmp/FY6900.pipe' \
--stdout=/tmp/PEMF-generator1.log \
/home/pi/git/PEMF/PEMF generator
