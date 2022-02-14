#! /bin/bash

/usr/bin/nice -n -4 /usr/bin/daemon -u pi:dialout -n frequency-serverd2 \
--env='HOME=/home/pi/git/frequency-server' \
--env='WEBPORT=8082' \
--env='GENFACTOR=1' \
--env='GENPORT=S' \
--env='PIPE=/tmp/FY6900.pipe' \
--stdout=/tmp/frequency-server2.log \
/home/pi/git/frequency-server/frequency-server
