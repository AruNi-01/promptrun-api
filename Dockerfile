FROM golang:latest

MAINTAINER AarynLu<hello@0x3f4.run>

WORKDIR /app/promptrun
COPY . .

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN go build -o promptrun .

EXPOSE 8080

ENTRYPOINT ["./promptrun"]
