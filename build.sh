# executable for windows.
GOOS=windows GOARCH=amd64 go build -o ./server-win.exe main.go config.go

# executable for linux.
GOOS=linux GOARCH=amd64 go build -o ./server-lin.exe main.go config.go