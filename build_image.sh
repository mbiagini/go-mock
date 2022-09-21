VERSION=1.0.0
APP_NAME=go-mock

docker rmi      localhost:5000/${APP_NAME}:${VERSION}
docker build -t localhost:5000/${APP_NAME}:${VERSION} .
docker push     localhost:5000/${APP_NAME}:${VERSION}