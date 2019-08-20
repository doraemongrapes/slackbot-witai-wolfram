FROM golang

ENV WOLFRAM_APP_ID="..."
ENV WITAI_ACCESS_TOKEN="..."
ENV SLACK_ACCESS_TOKEN="..."

# Add the source code:
ADD . /go/src/witai-wolfram-slackbot

# Build it:
RUN go get github.com/Krognol/go-wolfram
RUN go get github.com/christianrondeau/go-wit
RUN go get github.com/nlopes/slack

RUN go install witai-wolfram-slackbot

ENTRYPOINT /go/bin/witai-wolfram-slackbot

EXPOSE 8080
