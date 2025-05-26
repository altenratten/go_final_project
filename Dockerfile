#check=skip=SecretsUsedInArgOrEnv

#Команда запуска $ docker run -v $(pwd)/scheduler.db:/app/scheduler.db -e TODO_PORT=80 -p 80:80 todo

# Stage 1: Build the Go binary
FROM golang:1.24.2 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum 
COPY go.mod ./ 
COPY go.sum ./

#Downloads the dependencies
RUN go mod download

#Copy package and *.Go 
COPY pkg ./pkg
COPY *.go .

# Build the binary 
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 2
FROM ubuntu:latest

#Set ENV
ENV TODO_PORT=80 \
    TODO_DBFILE=/app/scheduler.db \
    TODO_PASSWORD=albatross6-send-married

# Copy the compiled binary and web
COPY web ./web
COPY --from=builder /app/server /usr/local/bin/server

# Expose port
EXPOSE 80

# Command to run the application
CMD ["/usr/local/bin/server"]