#!/bin/sh
echo $1 > flag
/usr/sbin/xinetd -dontfork