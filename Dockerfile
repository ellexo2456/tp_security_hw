FROM golang:1.21.1-alpine

COPY . app

RUN cd app && \
    go build -o server ./ && \
    chmod +x src/certs/gen.sh && \
    apk add openssl

EXPOSE 8080/tcp

CMD cd app && ./server
