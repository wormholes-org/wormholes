FROM ubuntu:20.04
RUN apt-get update
RUN apt-get install supervisor -y
WORKDIR /wm
COPY wormholes.conf /etc/supervisor/conf.d/
COPY wormholes /wm/
#COPY start.sh /wm/
RUN mkdir -p /wm/.wormholes/wormholes
CMD ["/usr/bin/supervisord", "-n"]

