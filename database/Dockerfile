FROM golang:1.24-alpine

WORKDIR database

# Install required tools: curl and ca-certificates
RUN apk add --no-cache curl ca-certificates

# Download migrate CLI
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


# Copy migration files into the image
COPY migrations migrations

# Add a shell script as entrypoint
COPY entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]