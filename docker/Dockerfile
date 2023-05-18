FROM alpine:3.14

RUN mkdir -p /usr/local/app/resources
RUN mkdir -p /usr/local/app/log
ADD zerver_lin.exe /usr/local/app/zerver_lin.exe

RUN apk --no-cache add curl

WORKDIR /usr/local/app
EXPOSE 5000