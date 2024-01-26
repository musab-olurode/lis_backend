# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Olurode Mus'ab <olurodemusab@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest

# Install Python, pip, make, go, and psycopg2
RUN apk add python3 py3-pip py3-psycopg2 make go

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 5000 to the outside
EXPOSE 5000

# Install goose
RUN go install github.com/pressly/goose/cmd/goose@latest

# Add the directory where goose is installed to the PATH
ENV PATH="/root/go/bin:${PATH}"

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Copy the 'entrypoint.sh' script into the Docker image
COPY entrypoint.sh /entrypoint.sh

# Make the 'entrypoint.sh' script executable
RUN chmod +x /entrypoint.sh

# Set the 'entrypoint.sh' script as the entrypoint
ENTRYPOINT ["/entrypoint.sh"]

# Command to run the executable
CMD ["./main"]