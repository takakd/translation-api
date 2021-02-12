FROM golang:1.14

WORKDIR /go/src/api
COPY ./src .

RUN go mod vendor
RUN make build

CMD ["cmd/api/api"]
