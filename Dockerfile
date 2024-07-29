FROM golang:1.22.4-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o req_3rd-party_services ./cmd/main.go

CMD ["./req_3rd-party_services"]