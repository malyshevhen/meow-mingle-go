# Builder
FROM golang:1.22.1-alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create appuser.
ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Set destination for COPY
WORKDIR $GOPATH/src/github.com/malyshEvhen/meow_mingle

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/meow_mingle cmd/api/main.go


# Runner
FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/meow_mingle /go/bin/meow_mingle
EXPOSE 8080

# Use an unprivileged user.
USER appuser:appuser

# Run the binary
ENTRYPOINT ["/go/bin/meow_mingle"]
