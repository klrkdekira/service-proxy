FROM ubuntu:14.04
MAINTAINER Chee Leong <klrkdekira@gmail.com>

RUN apt-get update -q && \
    apt-get upgrade -qy && \
    apt-get install -qy ca-certificates && \
    apt-get clean

ADD sherpa /usr/bin/sherpa

EXPOSE 8080

CMD ["/usr/bin/sherpa"]
