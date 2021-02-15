FROM golang:1.14

# Copy go codes to build.
WORKDIR /go/src/api
COPY ./src .

# Prepare modules.
RUN go mod vendor
# Build. For details about building, to see Makefile.
RUN make build

CMD ["cmd/api/api"]
