FROM golang:latest

WORKDIR /app/

COPY pocketbase.go .
COPY go.mod .
RUN go mod tidy
RUN go build -o pocketbase
EXPOSE 8090
RUN mv /app/pocketbase /usr/local/bin/pocketbase
ENTRYPOINT ["/usr/local/bin/pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/pb_data"]
