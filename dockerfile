# Use the official Golang image from the Docker Hub
FROM golang:1.22.6

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Get the dependencies
RUN go mod download

# Generate the gRPC code
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/v1/log.proto

# Build the Go application
RUN go build -o server ./server/main.go

# Make port available to the world outside this container
EXPOSE 8080

# Run the gRPC server
CMD ["./server"]