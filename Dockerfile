FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o test-loader

FROM scratch
WORKDIR /app
COPY --from=builder /app/test-loader .
ENTRYPOINT ["./test-loader"]