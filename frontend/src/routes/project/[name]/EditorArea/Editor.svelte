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
	import { json } from "@sveltejs/kit";


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
        byteCChanges64: string
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

            if (change.text.length > 1) {
                const newCursor = change.fromB + change.text.length;

                for (let i = 0; i < change.text.length; i++) {
                    
                    editorView.dispatch({
                    changes: {
                        from: change.fromB+i,
                        to: change.fromB+i,
                        insert: change.text.charAt(i),
                    },
                    selection: {
                        anchor: change.fromB+i+1
                    },
                });

                }

                
                return;
            }
            
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

        const fo: UpdatedDocMessage = UpdateDocument(changes, cursorIndex)

        const env: Envelope = {
            type: "operation",
            editDocMsg: fo
        }
        
        const json = JSON.stringify(env)
        socket.send(json)     

                
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

    }

    const BlockLocalChanges = EditorState.transactionFilter.of(tr => {
        if (tr.docChanged && !updateFromCode) {
            sendChangesToCrdt(tr);
        } 
        return tr
    })

    function onUpdate(update: ViewUpdate) {
        if (updateFromCode) return;
        return
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


    function decodeByteData(encodedState: string): string {

        const binaryStr = atob(encodedState);

        const bytes = Uint8Array.from(binaryStr, c => c.charCodeAt(0));

        const decoder = new TextDecoder("utf-8");
        const dataString = decoder.decode(bytes);

        return dataString
    }

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

        //socket = new WebSocket(`${serverUrl}/editDocWebsocket`);
        // in your front-end:
        socket = new WebSocket(`${serverUrl}/editDocWebsocket?projectName=${page.params.name}&documentName=document.tex`)

        socket.addEventListener("message", (event) => {

            //const text = new TextDecoder().decode(event.data);

            const message = JSON.parse(event.data);


            switch (message.type) {

                case "user_connected":
                    console.log("New user connected. ID: " + message.id);
                    SetUserID(message.id)

                    socket.send(JSON.stringify({
                        type:     "stateRequest",
                    }));
                    // todo: implementera nån wait function (promise?) och nån 
                    //      timeout om den inte får tillbaka CRDT state inom x sekunder

                    break;

                case "stateResponse":
                    updateFromCode = true;

                    const jsonString = decodeByteData(message.byteState as string)
                    
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

                    console.log(message.editDocMsg.byteCChanges);
                    applyUpdate(message.editDocMsg.byteCChanges)
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
