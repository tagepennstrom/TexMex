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
	"encoding/base64"
	"encoding/json"
	"syscall/js"

	"websocket-server/crdt"
)

func initUser(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		println("Wrong number of arguments [ SetUserID() ]")
		return nil
	}
	id := args[0].Int()

	crdt.SetUserID(id)

	return nil
}

func loadState(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		println("Wrong number of arguments [ InitializeDocument() ]")
		return nil
	}

	jsonString := args[0].String()

	newDocStr := crdt.LoadSnapshot(jsonString)

	return newDocStr
}

func updateDocumentWrap(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		println("Wrong number of arguments")
		return nil
	}

	changes := make([]crdt.Change, args[0].Length())
	for i := range len(changes) {
		change := args[0].Index(i)
		changes[i] = crdt.Change{
			FromA: change.Get("fromA").Int(),
			ToA:   change.Get("toA").Int(),
			FromB: change.Get("fromB").Int(),
			ToB:   change.Get("toB").Int(),
			Text:  change.Get("text").String(),
		}
	}

	cursorIndex := args[1].Int()

	res := crdt.UpdateDocument(changes, cursorIndex)

	jsonChanges, err := json.Marshal(res.CChanges)
	if err != nil {
		println("marshal coordChanges failed:", err)
		return nil
	}

	b64 := base64.StdEncoding.EncodeToString(jsonChanges)

	jsObj := js.Global().Get("Object").New()
	jsObj.Set("cursorIndex", res.CursorIndex)
	jsObj.Set("byteCChanges", b64)

	return jsObj
}

func handleOperation(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		println("Wrong number of arguments (HandleOperation)")
		return ""
	}

	b64Input := args[0].String() // b64 enkodat
	jsonBytes, err := base64.StdEncoding.DecodeString(b64Input)
	if err != nil {
		println("Failed to decode Base64: %v", err)
	}

	jsonString := string(jsonBytes)

	jsonIndexChanges := crdt.DocuMain.HandleCChange(jsonString)

	return jsonIndexChanges
}

func debugFeat(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		println("Wrong number of arguments [ CRDebug() ]")
		return nil
	}

	verbose := args[0].Bool()

	crdt.PrintDocument(verbose)
	return nil
}

func registerCallbacks() {
	js.Global().Set("UpdateDocument", js.FuncOf(updateDocumentWrap))
	js.Global().Set("SetUserID", js.FuncOf(initUser))
	js.Global().Set("CRDebug", js.FuncOf(debugFeat))
	js.Global().Set("LoadState", js.FuncOf(loadState))
	js.Global().Set("HandleOperation", js.FuncOf(handleOperation))

	println("Function callbacks registered")
}

func main() {
	println("WASM is alive")
	registerCallbacks()
	select {} // keep running
}
