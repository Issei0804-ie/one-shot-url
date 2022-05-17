FROM golang:1.18.1

WORKDIR /go/src

# timezone change
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

#RUN apt clean -y
#RUN apt update -y 
#RUN apt upgrade -y

RUN apt-get update && apt-get install -y wget

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN  curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

CMD dockerize -wait tcp://rdb:3306 `cd /go/src && go run main.go`
