cli:
	go build -mod vendor -o bin/emit cmd/emit/main.go
	go build -mod vendor -o bin/derive-media-json cmd/derive-media-json/main.go
