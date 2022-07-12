#!/bin/sh
echo $1 | tr -d '\n' | hexdump -ve '1/1 "%.2x\n"' > /home/HTTELab/data.txt