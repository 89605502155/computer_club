FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o computer-club ./cmd/main.go

CMD ["./computer-club"]