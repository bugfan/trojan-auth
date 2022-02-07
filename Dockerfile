FROM golang:1.17 AS build-env
MAINTAINER sean bugfan "908958194@qq.com"
ADD . /trojan-auth
WORKDIR /trojan-auth
RUN go build -o trojan-auth main.go

FROM alpine:3.13
    
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY --from=build-env /trojan-auth/trojan-auth /trojan-auth

RUN chmod +x /trojan-auth

ENTRYPOINT ["/trojan-auth"]

FROM alpine:3.13

WORKDIR /
RUN apk add --update --no-cache
RUN apk add --update vim && \
    apk add --update nano
    
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY --from=build-env /trojan-auth/trojan-auth /trojan-auth

RUN chmod +x /trojan-auth

CMD /trojan-auth -server -remote 
