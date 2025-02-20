FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
#ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY ./etc /app/etc
COPY ./static /app/static
RUN go build -ldflags="-s -w" -o /app/sonex_api_backend_go sonex_api_backend_go.go


#FROM scratch
FROM alpine:latest
# 安装运行需要的最小环境
RUN apk add --no-cache tzdata ca-certificates

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
# ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/sonex_api_backend_go /app/sonex_api_backend_go
COPY --from=builder /app/etc /app/etc
COPY --from=builder /app/static /app/static

CMD ["./sonex_api_backend_go", "-f", "etc/sonex_api_backend_go.yaml"]
