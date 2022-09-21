FROM alpine:3.14

RUN mkdir -p /usr/local/app/resources
ADD server-lin.exe /usr/local/app/server-lin.exe

RUN apk --no-cache add curl

WORKDIR /usr/local/app
EXPOSE 5000