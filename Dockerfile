FROM golang:1.24 as builder

WORKDIR /app
COPY . .
ENV GO111MODULE=on GOPROXY=https://goproxy.cn
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o . -v ./...

FROM scratch

WORKDIR /app
COPY --from=builder /app/depends-on /app/depends-on

ENTRYPOINT ["/app/depends-on"]
