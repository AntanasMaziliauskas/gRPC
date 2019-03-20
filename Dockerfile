
# Start from golang v1.11 base image
FROM golang:1.11 as builder

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . "$GOPATH/src/github.com/AntanasMaziliauskas/grpc"

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/AntanasMaziliauskas/grpc/cmd/node

# Download dependencies
RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/node .


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/node .

#EXPOSE 8080

CMD ["./node", "-storage", "mongo"] 
