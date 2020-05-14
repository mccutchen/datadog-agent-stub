# =============================================================================
# Build stage
# =============================================================================
FROM golang:1.14 AS build

ENV CGO_ENABLED=0
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make DIST_PATH=/bin

# =============================================================================
# Minimal final image
# =============================================================================
FROM gcr.io/distroless/base
COPY --from=build /bin/datadog-agent-stub /bin/datadog-agent-stub
CMD ["/bin/datadog-agent-stub"]
