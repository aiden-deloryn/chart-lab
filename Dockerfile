FROM golang:1.16 as builder

WORKDIR /go/src/

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o /go/bin/app ./pkg

FROM alpine:3.14.0

WORKDIR /chartlab

COPY --from=builder /go/bin/app /chartlab/app

COPY ssl/chartlab.crt /chartlab/ssl/

COPY ssl/chartlab.key /chartlab/ssl/

EXPOSE 80/tcp

EXPOSE 443/tcp

ENTRYPOINT [ "/chartlab/app" ]