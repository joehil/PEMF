#! /bin/bash

/usr/bin/daemon -u pi:dialout -n frequency-serverd2 \
--env='HOME=/home/pi/git/frequency-server' \
--env='WEBPORT=8082' \
--env='GENFACTOR=1000000' \
--stdout=/tmp/frequency-server2.log \
/home/pi/git/frequency-server/frequency-server
