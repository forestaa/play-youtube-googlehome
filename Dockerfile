FROM schachr/raspbian-stretch:latest

RUN apt-get update && apt-get install -y git gcc youtube-dl

ENV APP_ROOT /app
WORKDIR $APP_ROOT
ADD ./go1.11.5.linux-armv6l.tar.gz /usr/local/
COPY ./go.mod $APP_ROOT/
COPY ./src/go/server.go $APP_ROOT/
COPY ./views/template.html $APP_ROOT/views/
COPY ./views/index.html $APP_ROOT/views/

RUN /usr/local/go/bin/go build -o server || /usr/local/go/bin/go build -o server

CMD ["./server"]
