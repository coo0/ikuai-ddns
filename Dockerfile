FROM alpine:latest
 
WORKDIR /build
COPY ikuai-ddns .

RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

CMD ["./ikuai-ddns", "-c", "/etc/ikuai-ddns/config.yml"]