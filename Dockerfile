FROM golang:1.21

RUN go version
ENV GOPATH=/

RUN apt update && \
    apt install --no-install-recommends -y git wait-for-it

COPY ./ ./

RUN go mod download
RUN GOOS=linux go build -o ad-submission ./src/cmd/app

CMD ["./ad-submission"]
