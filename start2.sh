#! /bin/bash

/usr/bin/nice -n -4 /usr/bin/daemon -u pi:dialout -n PEMF-serverd2 \
--env='HOME=/home/pi/git/PEMF' \
--env='WEBPORT=8082' \
--env='GENFACTOR=1' \
--env='GENPORT=S' \
--env='PIPE=/tmp/FY6900.pipe' \
--stdout=/tmp/PEMF-server2.log \
/home/pi/git/PEMF/PEMF
