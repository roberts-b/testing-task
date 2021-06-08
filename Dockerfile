FROM golang:1.16.5-alpine3.13

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/spr
COPY . .

RUN apk update && apk add make && \
    rm -rf /var/cache/apk/*

RUN go get -d -v .
RUN make build
RUN ls

EXPOSE 8888

CMD ["./testing-task"]
