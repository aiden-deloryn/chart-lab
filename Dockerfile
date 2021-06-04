FROM golang:1.16

WORKDIR /go/src/

COPY . .

RUN go build -o /go/bin/app ./pkg

ENTRYPOINT [ "/go/bin/app" ]