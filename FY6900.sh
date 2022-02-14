#! /bin/bash

/usr/bin/nice -n -5 /usr/bin/daemon -u pi:dialout -n frequency-generatord1 \
--env='HOME=/home/pi/git/frequency-server' \
--env='USBSPEED=115200' \
--env='USBPORT=usb-1a86_USB_Serial-if00-port0' \
--env='PIPE=/tmp/FY6900.pipe' \
--stdout=/tmp/frequency-generator1.log \
/home/pi/git/frequency-server/frequency-server generator
