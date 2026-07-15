FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/gialosai ./src/main.go

FROM --platform=$TARGETPLATFORM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /out/gialosai /app/gialosai

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/gialosai"]