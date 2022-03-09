#! /bin/bash

/usr/bin/nice -n -4 /usr/bin/daemon -u pi:dialout -n PEMF-serverd1 \
--env='HOME=/home/pi/git/PEMF' \
--env='WEBPORT=8080' \
--env='GENFACTOR=1' \
--env='GENPORT=P' \
--env='PIPE=/tmp/FY6900.pipe' \
--stdout=/tmp/PEMF-server1.log \
/home/pi/git/PEMF/PEMF
