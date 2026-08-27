FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/gialosai ./cmd/gialosai

FROM --platform=$TARGETPLATFORM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /out/gialosai /app/gialosai

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/gialosai"]