# Use the official Golang image from the Docker Hub
FROM golang:1.22.6

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Get the dependencies
RUN go mod download

# Build the Go application
RUN go build -o gateway_bin gateway/gateway.go

# Make port available to the world outside this container
EXPOSE 8082

# Run the gRPC server
CMD ["./gateway_bin"]