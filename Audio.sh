#! /bin/bash

/usr/bin/nice -n -5 /usr/bin/daemon -u pi:dialout -n PEMF-audiod \
--env='HOME=/home/pi/git/PEMF' \
--env='PIPE=/tmp/Audio.pipe' \
--env='TONES=/usr/bin/tones' \
--stdout=/tmp/PEMF-audio.log \
/home/pi/git/PEMF/PEMF audio
