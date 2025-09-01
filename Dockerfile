FROM golang:1.24-bullseye AS builder

ENV PATH="/go/bin:${PATH}"
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build

COPY . .

# Build dependencies.
RUN make install

RUN ls -al /go/bin

# Build code.
RUN make build

# Build the final container.
# https://github.com/GoogleContainerTools/distroless/tree/main/base
FROM gcr.io/distroless/static

ENV PORT=8080

WORKDIR "/app"

# Add the binary
COPY --from=builder /build/bin .

CMD ["./rm-rf-production"]