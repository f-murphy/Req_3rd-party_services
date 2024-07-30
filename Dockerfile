FROM golang:alpine

RUN go version
ENV GOPATH=/

COPY ./ ./
COPY ./configs /configs

RUN go mod download
RUN go build -o Req_3rd-party_services ./cmd/main.go

CMD ["./Req_3rd-party_services"]