FROM golang:1.24-bullseye AS builder

WORKDIR /build

COPY . .

ENV GOBIN=/usr/local/bin
ENV PATH="/usr/local/bin:${PATH}"

# Build dependencies.
RUN make install

ENV CGO_ENABLED=0

# Build code.
RUN make build

# Avoid breaking functionality by running tests
RUN make test

# Build the final container.
# https://github.com/GoogleContainerTools/distroless/tree/main/base
FROM gcr.io/distroless/static

ENV PORT=8080

WORKDIR "/app"

# Add the binary
COPY --from=builder /build/bin .

CMD ["./rm-rf-production"]