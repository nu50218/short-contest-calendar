FROM golang:1.18

# Set the Current Working Directory inside the container
WORKDIR /app
COPY . .
RUN go build -o ./app main.go

ENTRYPOINT ["./app"]
