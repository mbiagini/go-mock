# executable for windows.
GOOS=windows GOARCH=amd64 go build -o ./zerver_win.exe main.go config.go router.go

# executable for linux.
GOOS=linux GOARCH=amd64 go build -o ./zerver_lin.exe main.go config.go router.go