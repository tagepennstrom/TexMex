<script lang='ts'>
    import {basicSetup, EditorView} from "codemirror"
    import { page } from '$app/state'
    import {onMount} from 'svelte'
    import {EditorState, Transaction} from "@codemirror/state"
    import {StreamLanguage,} from '@codemirror/language'
    import { stex } from "@codemirror/legacy-modes/mode/stex"
    import { editorView as editorViewStore } from "$lib/stores";
    import { autocompletion } from "@codemirror/autocomplete";
    import { myCompletions } from "$lib/completions";


    let socket: WebSocket;

    let updateFromCode = $state(false);

    let editor: HTMLElement;
    let editorView: EditorView;

    let uID: number; // saved uID

    type Change = {
        fromA: number;   // Start index
        toA: number;     // Slut index
        fromB: number;   // Start index
        toB: number;     // Slut index
        text: string;   // Tillagd text, tom vid borttagning
    }

    type UpdatedDocMessage = {
        document: string
        cursorIndex: number
        jsonCChanges: string
    }

    type Envelope = {
        Type: string          
        EditDocMsg: UpdatedDocMessage
    }

    async function applyUpdate(document: string, jsonChanges: string) {
        updateFromCode = true;

        const strDoc: string = HandleOperation(jsonChanges)

        // TODO: Hårdkodad cursor index for now
        const cursorIndex: number = strDoc.length

        editorView.dispatch({
            changes: {
                from: 0,
                to: editorView.state.doc.length,
                insert: strDoc,
            },
            selection: {
                anchor: cursorIndex,
            },
        });
        updateFromCode = false;
    }


    function sendChangesToCrdt(tr: Transaction): void {
        
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

        const document = editorView.state.doc.toString();
        const cursorIndex = editorView.state.selection.main.anchor;

        const updatedDocMessage: UpdatedDocMessage = UpdateDocument(
            document,
            changes,
            cursorIndex,
        )

        const env: Envelope = {
            Type: "operation",
            EditDocMsg: updatedDocMessage,

        }

        console.log("Sending envelope:",env)

        socket.send(JSON.stringify(env));

        const serverUrl = `http://${location.hostname}:8080`;
        fetch(`${serverUrl}/projects/${page.params.name}/documents/document.tex`, {
            method: "PUT",
            headers: { "Content-Type": "text/plain" },
            body: updatedDocMessage.document,
        })
    }




    const BlockLocalChanges = EditorState.transactionFilter.of(tr => {
        if (tr.docChanged && !updateFromCode) {
            sendChangesToCrdt(tr);
        }
        return tr;
    })

    
    const fixedHeightEditor = EditorView.theme({
        "&": {height: "700px"},
        ".cm-scroller": {overflow: "auto"}
    })

    const extensions = [
        basicSetup,
        StreamLanguage.define(stex),
        fixedHeightEditor,
        BlockLocalChanges,
        EditorView.lineWrapping,
        autocompletion({ override: [myCompletions] })]



    onMount(() => {
        const serverUrl = `http://${location.hostname}:8080`;
        socket = new WebSocket(`${serverUrl}/editDocWebsocket`);

        socket.addEventListener("message", (event) => {
            const message = JSON.parse(event.data);
            console.log("hello this is type:", message.type)
            switch (message.type) {

                case "user_connected":
                    console.log("New user connected. ID: " + message.id);
                    SetUserID(message.id)
                    uID = message.id

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

                    console.log("Recieved doc:", loadedDocStr)
                    console.log("Dispatched text to frontend.")

                    const cursorPos = loadedDocStr.length // onödigt?

                    editorView.dispatch({
                                changes: {
                                    from: 0,
                                    to: editorView.state.doc.length,
                                    insert: loadedDocStr,
                                },
                                selection: {
                                    anchor: cursorPos
                                },
                            });

                    updateFromCode = false;
                    
                    break;
                    

                case "operation":
                    console.log(message.editDocMsg);

                    let changes = message.editDocMsg.jsonCChanges
                    applyUpdate(editorView.state.doc.toString(), changes)
                    break;


                default:
                    console.log("Error: No switch-case handler found for message.");
                    console.log(message)

                    break;

            }


            if (message.id != undefined){
                

            } else {
                
            }

        });

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
