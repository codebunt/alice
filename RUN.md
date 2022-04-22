go build -buildmode=c-shared -o ./bin/mpc.so lib/mpc.go


GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o /Users/shetty/code/kapow/kapow-mpc/dart-mpc/lib/src/native/mpc.dylib lib/mpc.go

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o /Users/shetty/code/kapow/kapow-mpc/dart-mpc/lib/src/native/mpc.so lib/mpc.go

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o /Users/shetty/code/kapow/kapow-mpc/dart-mpc/lib/src/native/mpc.a lib/mpc.go


https://www.digitalocean.com/community/tutorials/building-go-applications-for-different-operating-systems-and-architectures


GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o /Users/shetty/code/kapow/kapow-mpc/dart-mpc/lib/src/native/mpc.so lib/mpc.go


GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o /Users/shetty/code/kapow/kapow-mpc/dart-mpc/lib/src/native/mpc.so lib/mpc.go

/root/kapow-mpc/java-server
