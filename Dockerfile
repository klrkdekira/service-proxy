FROM ubuntu:14.04
MAINTAINER Chee Leong <klrkdekira@gmail.com>

RUN apt-get update -q && \
    apt-get upgrade -qy && \
    apt-get clean

ADD go-api-mirror /usr/bin/go-api-mirror

EXPOSE 8080

CMD ["/usr/bin/go-api-mirror"]
