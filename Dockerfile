FROM golang:latest
ARG bot_id
ADD . /go/src/github.com/rpeshkov/slovnik_bot
WORKDIR /go/src/github.com/rpeshkov/slovnik_bot
RUN go get github.com/tools/godep
RUN godep restore
RUN go install github.com/rpeshkov/slovnik_bot
ENV SLOVNIK_BOT_ID $bot_id
ENTRYPOINT /go/bin/slovnik_bot
