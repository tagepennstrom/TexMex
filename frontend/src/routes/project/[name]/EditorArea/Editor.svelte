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
        HandleOperation(jsonChanges);
        const changes: CoordChanges[] = JSON.parse(jsonChanges);
        updateFromCode = true;
        changes.forEach((change) => {
            const index = CoordinateToIndex(JSON.stringify(change.coordinate));
            editorView.dispatch({
                changes: {
                    from: index,
                    to: index + change.letter.length,
                    insert: change.letter,
                },
                selection: {
                    anchor: GetCursorIndex(),
                },
            });
        })
        updateFromCode = false;
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

        const actualChanges: TransactionSpecChange[] = coordChanges.map(change => {
            const index = CoordinateToIndex(JSON.stringify(change.coordinate));
            if (change.operation === "delete") {
                return {
                    from: index,
                    to: index + 1,
                    insert: change.letter,
                };
            } else {
                return {
                    from: index,
                    to: index,
                    insert: change.letter,
                };
            }
        });
        return {
            changes: actualChanges,
            selection: {
                anchor: updDocMsg.cursorIndex,
            }
        };
    }

    const BlockLocalChanges = EditorState.transactionFilter.of(tr => {
        if (tr.docChanged && !updateFromCode) {
            return sendChangesToCrdt(tr);
        } else {
            return tr;
        }
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
