FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH BUILDPLATFORM
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o netviz .

FROM alpine:3.23

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/netviz /netviz

EXPOSE 8080

ENTRYPOINT ["/netviz"]
