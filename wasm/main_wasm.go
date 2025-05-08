/**

För att kompilera:
1. Gå till /wasm directory
2. Kör: GOOS=js GOARCH=wasm go build -o ../frontend/src/wasm/main.wasm
2.1 powershell: $env:GOOS="js"; $env:GOARCH="wasm"; go build -o ..\frontend\src\wasm\main.wasm
2.2 cmd.exe:
set GOOS=js
set GOARCH=wasm
go build -o ..\frontend\src\wasm\main.wasm

**/

package main

import (
	"syscall/js"

	"websocket-server/crdt"
)

func initUser(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		println("Wrong number of arguments")
		return nil
	}
	id := args[0].Int()
	println("Initializing user ", id, " (main_wasm)")

	crdt.SetUserID(id)

	return nil
}

func updateDocumentWrap(this js.Value, args []js.Value) any {
	if len(args) != 3 {
		println("Wrong number of arguments")
		return nil
	}

	document := args[0].String()
	cursorIndex := args[2].Int()

	changes := make([]crdt.Change, args[1].Length())
	for i := range len(changes) {
		change := args[1].Index(i)
		changes[i] = crdt.Change{
			FromA: change.Get("fromA").Int(),
			ToA:   change.Get("toA").Int(),
			FromB: change.Get("fromB").Int(),
			ToB:   change.Get("toB").Int(),
			Text:  change.Get("text").String(),
		}
	}
	res := crdt.UpdateDocument(document, changes, cursorIndex)
	var m = make(map[string]any)
	m["document"] = res.Document
	m["cursorIndex"] = res.CursorIndex

	return js.ValueOf(m)
}

func registerCallbacks() {
	js.Global().Set("UpdateDocument", js.FuncOf(updateDocumentWrap))
	js.Global().Set("SetUserID", js.FuncOf(initUser))

	println("Function callbacks registered")
}

func main() {
	println("WASM is alive")
	registerCallbacks()
	select {} // keep running
}
