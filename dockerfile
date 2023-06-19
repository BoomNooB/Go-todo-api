FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY local.env ./

RUN go mod download

COPY . ./

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

RUN go build \
    -ldflags "-X main.buildCommit=`git rev-parse --short HEAD` \
    -X main.buildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" \
    -o /main/app

# Deploy

FROM gcr.io/distroless/base-debian11

COPY --from=builder /main/app /main/app

EXPOSE 8081

USER nonroot:nonroot

CMD ["/main/app"]
