# Use a lightweight base image
FROM alpine:latest

# Install Go
RUN apk add --no-cache go

# Set the working directory
WORKDIR /app

# Copy the Go application code
COPY . .

# Build the Go application
RUN go build -o myapp .

# Command to run the application
CMD ["./myapp"]
