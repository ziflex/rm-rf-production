FROM alpine:3.22

WORKDIR database

# Install required tools: curl and ca-certificates
RUN apk add --no-cache curl ca-certificates

# Download migrate CLI
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz \
    | tar xvz -C /usr/local/bin/ \
    && chmod +x /usr/local/bin/migrate

# Copy migration files into the image
COPY migrations migrations

# Add a shell script as entrypoint
COPY entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]