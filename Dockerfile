FROM golang:1.24-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

RUN echo "Building for OS=$TARGETOS ARCH=$TARGETARCH VERSION=$VERSION COMMIT=$COMMIT BUILD_DATE=$BUILD_DATE"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/app ./cmd/app
COPY internal ./internal


RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath \
    -ldflags "-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.buildDate=$BUILD_DATE" \
    -o backend ./cmd/app

FROM alpine:3.19 AS runner
WORKDIR /app

COPY --from=builder /app/backend ./backend
COPY ./configs /configs

ENTRYPOINT ["./backend", "--config", "/configs/local.yaml"]