FROM alpine:latest
 
WORKDIR /build
COPY ikuai-ddns .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

CMD ["./ikuai-ddns", "-c", "/etc/ikuai-ddns/config.yml"]