#!/bin/sh
echo $1 | tr -d '\n' | hexdump -ve '1/1 "%.2x\n"' > /home/HTTELab/data.txt
chmod +x /home/HTTELab/server
cd /home/HTTELab/ && ./server --port 8888 --capture traces.json.gz -static "traces.json.gz" -static "elmo.tgz" -static "firmware.tgz"
