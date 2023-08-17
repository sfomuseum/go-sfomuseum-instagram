GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/emit cmd/emit/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/derive-media-json cmd/derive-media-json/main.go

training-data:
	curl -o data/english.json https://raw.githubusercontent.com/neurosnap/sentences/main/data/english.json
