
FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /go-service-boilerplate

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /go-service-boilerplate/main .

# Build a small image
FROM scratch

COPY --from=builder /dist/main /
COPY env/APP_ENV.env /app

EXPOSE 8080
# Command to run
ENTRYPOINT ["/main"]