#!/bin/sh
# Add your startup script

# DO NOT DELETE
echo $1 > flag
export PORT=`cat port | awk -F: '{print $4}'`
./pfs -connect $PORT -serve 8888
/usr/sbin/xinetd -dontfork
