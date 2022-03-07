#! /bin/bash

/usr/bin/nice -n -5 /usr/bin/daemon -u pi:dialout -n frequency-audiod \
--env='HOME=/home/pi/git/frequency-server' \
--env='PIPE=/tmp/Audio.pipe' \
--env='TONES=/usr/bin/tones' \
--stdout=/tmp/frequency-audio.log \
/home/pi/git/frequency-server/frequency-server audio
