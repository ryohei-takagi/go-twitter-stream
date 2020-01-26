FROM alpine:3.11.2

RUN mkdir /next
WORKDIR /next

RUN apk --update add tzdata nodejs npm && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*
