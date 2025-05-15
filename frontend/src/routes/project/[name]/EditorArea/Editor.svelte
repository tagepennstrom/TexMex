<script lang='ts'>
    import {basicSetup, EditorView} from "codemirror"
    import { page } from '$app/state'
    import {onMount} from 'svelte'
    import {EditorState, Transaction, type TransactionSpec } from "@codemirror/state"
    import { ViewUpdate } from "@codemirror/view"
    import {StreamLanguage,} from '@codemirror/language'
    import { stex } from "@codemirror/legacy-modes/mode/stex"
    import { editorView as editorViewStore } from "$lib/stores";
    import { autocompletion } from "@codemirror/autocomplete";
    import { myCompletions } from "$lib/completions";


    let socket: WebSocket;

    let updateFromCode = $state(false);

    let editor: HTMLElement;
    let editorView: EditorView;

    type Change = {
        fromA: number;   // Start index
        toA: number;     // Slut index
        fromB: number;   // Start index
        toB: number;     // Slut index
        text: string;   // Tillagd text, tom vid borttagning
    }

    type CoordT = {
        coordinate: number[];
        id: number;
    }

    type CoordChanges = {
        coordinate: CoordT;
        operation: string;
        letter: string;
    }

    type UpdatedDocMessage = {
        cursorIndex: number
        jsonCChanges: string
    }

    type Envelope = {
        type: string          
        editDocMsg: UpdatedDocMessage
    }

    async function applyUpdate(jsonChanges: string) {
        const jsonIndexChanges: string = HandleOperation(jsonChanges);
        const iChanges: Change[] = JSON.parse(jsonIndexChanges)

        updateFromCode = true;

        iChanges.forEach((change) => {
            const cursorPos = editorView.state.selection.main.anchor;

            let newCursor: number = 0;

            if (change.text == "") {
                newCursor = cursorPos > change.fromB ? cursorPos - 1 : cursorPos
                
            } else {
                newCursor = cursorPos > change.fromB ? cursorPos + 1 : cursorPos
            }

            editorView.dispatch({
                    changes: {
                        from: change.fromB,
                        to: change.toB,
                        insert: change.text,
                    },
                    selection: {
                        anchor: newCursor
                    },
                });
        })
        updateFromCode = false;
    }

    function isInsert(coordChanges: CoordChanges[]) {
        for (const change of coordChanges) {
            if (change.operation !== "insert") {
                return false;
            }
        }
        return true;
    }

    function isDelete(coordChanges: CoordChanges[]) {
        for (const change of coordChanges) {
            if (change.operation !== "delete") {
                return false;
            }
        }
        return true;
    }

    function totalDeletes(coordChanges: CoordChanges[]) {
        let i = 0;
        for (const change of coordChanges) {
            if (change.operation == "delete") {
                i++
            }
        }
        return i;
    }


    function sendChangesToCrdt(tr: Transaction): TransactionSpec {
        
        const changes: Change[] = [];
        tr.changes.iterChanges((fromA, toA, fromB, toB, inserted) => {
            changes.push({
                fromA: fromA, 
                toA: toA,     
                fromB: fromB, 
                toB: toB,     
                text: inserted.toString() // Tillagd text, tom vid borttagning
            });
        });

        const cursorIndex = editorView.state.selection.main.anchor;

        const updDocMsg: UpdatedDocMessage = UpdateDocument(
            changes,
            cursorIndex,
        )


        const env: Envelope = {
            type: "operation",
            editDocMsg: updDocMsg,
        }
        console.log("Sending envelope:",env)
        socket.send(JSON.stringify(env));

        
        const coordChanges: CoordChanges[] = JSON.parse(updDocMsg.jsonCChanges);
        
        type TransactionSpecChange = {
            from: number;
            to: number;
            insert: string;
        };

        const s: TransactionSpecChange = {
            from: 1,
            to:1,
            insert:""
        }

        return {
            changes: s,
            selection: {
                anchor: 0,
            }
        };
        

        let newCursorPos: number = updDocMsg.cursorIndex

        const actualChanges: TransactionSpecChange[] = [];
        if (isInsert(coordChanges)) {
            for (let i = 0; i < coordChanges.length; i++) {
                const change = coordChanges[i];
                const index = CoordinateToIndex(JSON.stringify(change.coordinate))-1;
                actualChanges.push({
                    from: index - i,
                    to: index - i,
                    insert: change.letter,
                });
            }
        } else if (isDelete(coordChanges)) {
            for (let i = 0; i < coordChanges.length; i++) {
                const change = coordChanges[i];
                const index = CoordinateToIndex(JSON.stringify(change.coordinate)) - 1;
                actualChanges.push({
                    from: index + i,
                    to: index + i + 1,
                    insert: "",
                });
            }
        } else {
            // SELECT AND REPLACE Operation
            // TODO: fungerar inte

            let tot: number = totalDeletes(coordChanges)
            
            CRDebug(true)

            for (let i = 0; i < tot; i++) {
                const change = coordChanges[i];
                let index = CoordinateToIndex(JSON.stringify(change.coordinate)) -2;
                if (index == -1) {
                    index = 0;
                }

                console.log("sel and del - i",i,"ind:",index,"pos", index + i, "coord:", change.coordinate.coord)
                actualChanges.push({
                    from: index + i,
                    to: index + i + 1,
                    insert: "",
                });
            }

            for (let i = 0; i < coordChanges.length - tot; i++) {
                const change = coordChanges[i + tot];
                
                let index = CoordinateToIndex(JSON.stringify(change.coordinate))-1;

                console.log("sel and ins:",i,"at:",index, "is:", change.letter)
                continue
                actualChanges.push({
                    from: index - i,
                    to: index - i,
                    insert: change.letter,
                });
            }
            
            for (let i = 0; i < actualChanges.length; i++){
                console.log(actualChanges[i])
            }
            newCursorPos = 0
            
        }

        return {
            changes: actualChanges,
            selection: {
                anchor: newCursorPos,
            }
        };
    }

    const BlockLocalChanges = EditorState.transactionFilter.of(tr => {
        if (tr.docChanged && !updateFromCode) {
            sendChangesToCrdt(tr);
        } 
        return tr
    })

    function onUpdate(update: ViewUpdate) {
        if (updateFromCode) return;
        const serverUrl = `http://${location.hostname}:8080`;
        fetch(`${serverUrl}/projects/${page.params.name}/documents/document.tex`, {
            method: "PUT",
            headers: { "Content-Type": "text/plain" },
            body: update.state.doc.toString(),
        })
    }

    const fixedHeightEditor = EditorView.theme({
        "&": {height: "700px"},
        ".cm-scroller": {overflow: "auto"}
    })

    const extensions = [
        basicSetup,
        StreamLanguage.define(stex),
        fixedHeightEditor,
        BlockLocalChanges,
        EditorView.updateListener.of(onUpdate),
        EditorView.lineWrapping,
        autocompletion({ override: [myCompletions] })
    ]


    onMount(() => {
        const serverUrl = `http://${location.hostname}:8080`;

        fetch(`${serverUrl}/projects/${page.params.name}/documents/document.tex`)
            .then(res => res.text())
            .then(text => {
                // Initialize CodeMirror editor
                editorView = new EditorView({
                    state: EditorState.create({
                        doc: text,
                        extensions
                    }),
                    parent: editor
                });

                // Spara editorn i store för delning med Toolbar
                editorViewStore.set(editorView);
            });

        socket = new WebSocket(`${serverUrl}/editDocWebsocket`);
        socket.addEventListener("message", (event) => {
            const message = JSON.parse(event.data);
            switch (message.type) {

                case "user_connected":
                    console.log("New user connected. ID: " + message.id);
                    SetUserID(message.id)

                    socket.send(JSON.stringify({
                        type:     "stateRequest",
                        targetId: message.id
                    }));
                    // todo: implementera nån wait function (promise?) och nån 
                    //      timeout om den inte får tillbaka CRDT state inom x sekunder

                    break;

                case "stateResponse":
                    updateFromCode = true;
                    const encodedState = message.byteState as string;

                    const jsonString = atob(encodedState);
                    const loadedDocStr = LoadState(jsonString)

                    console.log("Recieved doc state:", loadedDocStr)

                    editorView.dispatch({
                                changes: {
                                    from: 0,
                                    to: editorView.state.doc.length,
                                    insert: loadedDocStr,
                                },
                                selection: {
                                    anchor: loadedDocStr.length
                                },
                            });

                    updateFromCode = false;
                    
                    break;
                    

                case "operation":
                    console.log(message.editDocMsg);
                    applyUpdate(message.editDocMsg.jsonCChanges)
                    break;


                default:
                    console.log("Error: No switch-case handler found for message.");
                    console.log(message)

                    break;
            }

        });
    });
</script>

<div id="editor" bind:this={editor}></div>

<style>
    #editor {
        height: 700px;
        width: 49%;
        float: left;
        margin: auto;
    }

</style>
