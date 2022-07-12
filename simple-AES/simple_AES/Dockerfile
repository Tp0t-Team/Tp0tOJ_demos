FROM ubuntu:20.04
RUN sed -i "s/http:\/\/archive.ubuntu.com/http:\/\/mirrors.ustc.edu.cn/g" /etc/apt/sources.list
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
      apt-utils \
      build-essential \
      git \
      bzip2 \
      # xinetd \
      bsdmainutils \
      python3 \
      python3-pip \
      bsdmainutils \
      wget && \
      pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple && \
      python3 -m pip install pycryptodome numpy scipy && \
    apt-get clean
RUN useradd -u 8888 -m HTTELab

# Key is the first 16 bytes in data.txt

# ELMO executable and runtime dependencies.
COPY share/elmo.tgz /home/HTTELab/
RUN cd /home/HTTELab && tar xfvz elmo.tgz --strip=1

# Target firmware binary.
COPY share/firmware.tgz /home/HTTELab/firmware.tgz
RUN cd /home/HTTELab && tar xfvz firmware.tgz --strip=1 firmware/firmware.bin

# server
COPY share/server.tgz /home/HTTELab/server.tgz
RUN cd /home/HTTELab && tar xfvz server.tgz --strip=1

# Collect plaintext (randdata), ciphertext (printdata) and traces.
COPY share/traces/collect.py /home/HTTELab/traces/collect.py
COPY share/traces/downsample.py /home/HTTELab/traces/downsample.py

COPY entrypoint.sh /home/HTTELab/entrypoint.sh
ENTRYPOINT ["/home/HTTELab/entrypoint.sh"]

# RUN rm /etc/xinetd.d/*
# COPY xinetd /etc/xinetd.d/xinetd

EXPOSE 8888
