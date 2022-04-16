go build -buildmode=c-shared -o bin/mpc.so lib/mpc.go


GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o /Users/shetty/code/kapow/kapow-mpc/dart-mpc/lib/src/native/mpc.dylib lib/mpc.go

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o /Users/shetty/code/kapow/kapow-mpc/dart-mpc/lib/src/native/mpc.so lib/mpc.go

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o /Users/shetty/code/kapow/kapow-mpc/dart-mpc/lib/src/native/mpc.a lib/mpc.go