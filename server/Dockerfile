# 构建阶段
FROM golang:alpine AS builder


WORKDIR /go/src/github.com/flipped-aurora/gin-vue-admin/server
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go mod tidy \
    && go build -o server .

# 运行阶段
FROM alpine:latest

LABEL MAINTAINER="SliverHorn@sliver_horn@qq.com"

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk update && apk add --no-cache tzdata openntpd \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /go/src/github.com/flipped-aurora/gin-vue-admin/server

# 复制构建阶段的文件
COPY --from=builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/server ./
COPY --from=builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/resource ./resource/
COPY --from=builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/config.docker.yaml ./ 

EXPOSE 8888

# 使用 JSON 格式指定 ENTRYPOINT
ENTRYPOINT ["./server", "-c", "config.docker.yaml"]
