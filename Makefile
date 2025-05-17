.PHONY: wasm
.PHONY: client


server:
	clear; clear; \
	go run server.go project.go


client:
	clear; \
	cd ./wasm && GOOS=js GOARCH=wasm go build -o ../frontend/src/wasm/main.wasm
	cd frontend && npm run dev

wasm:
	cd ./wasm && GOOS=js GOARCH=wasm go build -o ../frontend/src/wasm/main.wasm
