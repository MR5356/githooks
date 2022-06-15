FROM golang:alpine AS builder
ENV GOPROXY="https://goproxy.cn,direct"
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go build -o githooks .

FROM alpine
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /app/githooks .
CMD ["./githooks"]

