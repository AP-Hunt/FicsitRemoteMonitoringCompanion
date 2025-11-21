FROM node:16-alpine AS map-builder
WORKDIR /build
COPY map/package*.json ./
RUN npm install
COPY map/ ./
RUN npm run compile

FROM golang:1.25.4-alpine AS go-builder
WORKDIR /build
COPY Companion/go.* ./
RUN go mod download
COPY Companion/ ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${VERSION}'" -o companion main.go

FROM alpine:latest
ARG PROMETHEUS_VERSION=2.27.1
ARG TARGETARCH
WORKDIR /app

RUN apk --no-cache add curl tar && \
    curl -sL https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-${TARGETARCH}.tar.gz | tar xz && \
    mv prometheus-${PROMETHEUS_VERSION}.linux-${TARGETARCH} prometheus && \
    apk del curl tar

COPY --from=go-builder /build/companion /app/
COPY --from=map-builder /build/index.html /build/map-16k.png /build/vendor /build/img /build/js /app/map/
COPY Companion/prometheus.yml /app/prometheus/prometheus.yml

ENV FRM_LOG_STDOUT=true
EXPOSE 9000

CMD ["/app/companion"]
