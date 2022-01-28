#! /bin/bash

/usr/bin/daemon -u pi:dialout -n frequency-serverd --env='HOME=/home/pi/git/frequency-server' --stdout=/tmp/frequency-server.log /home/pi/git/frequency-server/frequency-server
