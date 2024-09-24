FROM golang:1.12-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/ldh-dns

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
# RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN cd pkg && go build -o /tmp/ldh-dns/out/ldh-dns .

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/ldh-dns/out/ldh-dns /app/ldh-dns

# This container exposes port 8053 to the outside world
EXPOSE 8053

# Run the binary program produced by `go install`
CMD ["/app/ldh-dns"]