BASE_REPOSITORY=localhost:5000
APP_NAME=go-mock
VERSION=2.0.0

cd ..
./server_compile.sh
cp zerver_lin.exe ./docker
cd ./docker

docker rmi   -f ${BASE_REPOSITORY}/${APP_NAME}:${VERSION}
docker build -t ${BASE_REPOSITORY}/${APP_NAME}:${VERSION} .

rm zerver_lin.exe