FROM ubuntu:20.04
RUN sed -i "s/http:\/\/archive.ubuntu.com/http:\/\/mirrors.ustc.edu.cn/g" /etc/apt/sources.list
RUN apt-get update
RUN apt-get -y upgrade
RUN apt-get install -y apt-utils lib32z1 xinetd
RUN useradd -u 8888 -m pwn
COPY share/libunicorn.so.1 /usr/local/lib/libunicorn.so.1
RUN chmod 755 /usr/local/lib/libunicorn.so.1
RUN ldconfig
COPY share/easiestpwn /home/pwn/easiestpwn
RUN chmod 755 /home/pwn/easiestpwn
RUN rm /etc/xinetd.d/*
COPY xinetd /etc/xinetd.d/xinetd
COPY entrypoint.sh /home/pwn/entrypoint.sh
ENTRYPOINT ["/home/pwn/entrypoint.sh"]
# CMD ["/usr/sbin/xinetd", "-dontfork"]
EXPOSE 8888
