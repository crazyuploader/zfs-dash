FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH BUILDPLATFORM
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o zfs-dash .

FROM alpine:3.23

RUN apk --no-cache add ca-certificates tzdata wget

COPY --from=builder /app/zfs-dash /zfs-dash

EXPOSE 8080

ENTRYPOINT ["/zfs-dash"]
CMD ["serve"]
