FROM ubuntu:18.04
RUN sed -i "s/http:\/\/archive.ubuntu.com/http:\/\/mirrors.tuna.tsinghua.edu.cn/g" /etc/apt/sources.list && \
         sed -i "s/http:\/\/security.ubuntu.com/http:\/\/mirrors.tuna.tsinghua.edu.cn/g" /etc/apt/sources.list


RUN apt-get update && \
       apt-get -y dist-upgrade && \
		apt-get install -y lib32z1 xinetd 


RUN useradd -u 8888 -m ctf && \
	mkdir -p /var/www


COPY ./html.tar.gz /var/www/html.tar.gz
COPY ./libCoroutine.so /usr/lib/libCoroutine.so
RUN chmod 755 /usr/lib/libCoroutine.so
RUN ldconfig
COPY ./rhttpd /home/ctf/rhttpd
RUN rm /etc/xinetd.d/*
COPY ./xinetd /etc/xinetd.d/xinetd
COPY ./pfs /home/ctf/pfs
COPY ./entrypoint.sh /home/ctf/entrypoint.sh

RUN chmod 755 /home/ctf/pfs

RUN tar -xzvf /var/www/html.tar.gz -C /var/www/ && \
	chown -R root:ctf /home/ctf && \
	chmod -R 755 /home/ctf && \
	echo "Blocked by ctf_xinetd" > /etc/banner_fail && \
	chmod +x /home/ctf/entrypoint.sh

ENTRYPOINT ["/home/ctf/entrypoint.sh"]
EXPOSE 8888
