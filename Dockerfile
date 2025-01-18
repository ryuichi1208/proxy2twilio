# Use an official Go image as the base image
FROM golang:1.20 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifest and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o proxy-server main.go

# Use a lightweight image for the final build
FROM debian:bullseye-slim

# Set up a non-root user
RUN useradd -m proxy-user

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/proxy-server /app/proxy-server

# Set permissions and switch to the non-root user
RUN chown -R proxy-user:proxy-user /app
USER proxy-user

# Expose the port the application runs on
EXPOSE 3000

# Command to run the application
CMD ["./proxy-server"]
