# Use the official Golang image from the Docker Hub
FROM golang:1.19-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Compile protobuf files
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/v1/log.proto

# Build the Go application
RUN go mod tidy
RUN go build -o server ./server/main.go

# Make port 50051 available to the world outside this container
EXPOSE 8080

# Run the gRPC server
CMD ["./server"]