FROM golang:1.15.8 as build-env

WORKDIR /go/src/fennec
ADD . /go/src/fennec

RUN go get -d -v ./...

RUN go build -o /go/bin/fennec

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/fennec /

CMD ["/fennec", "/assets/config.json"]