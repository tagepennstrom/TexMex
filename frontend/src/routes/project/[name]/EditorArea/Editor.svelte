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

    type UpdatedDocMessage = {
        cursorIndex: number
        byteCChanges64: string
    }

    type Envelope = {
        type: string          
        editDocMsg: UpdatedDocMessage
    }

    function pasteUpdate(change: Change): number{
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

            return newCursor;
        }
        else return 0
    }

    async function applyUpdate(jsonChanges: string) {
        const jsonIndexChanges: string = HandleOperation(jsonChanges);
        const iChanges: Change[] = JSON.parse(jsonIndexChanges)

        updateFromCode = true;

        iChanges.forEach((change) => {
            const cursorPos = editorView.state.selection.main.anchor;

            let newCursor = pasteUpdate(change) // hanterar om change är en paste


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


    function sendChangesToCrdt(tr: Transaction){
        
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
        console.log("Operation sent")


    }

    const BlockLocalChanges = EditorState.transactionFilter.of(tr => {
        if (tr.docChanged && !updateFromCode) {
            sendChangesToCrdt(tr);
        } 
        return tr
    })

    const fixedHeightEditor = EditorView.theme({
        "&": {height: "700px"},
        ".cm-scroller": {overflow: "auto"}
    })

    function onUpdate(){
        return // Den här funktionen används inte längre, bara placeholder
    }

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

    function dispatchStateToEditor(loadedDocStr: string){
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

        socket = new WebSocket(`${serverUrl}/editDocWebsocket?projectName=${page.params.name}&documentName=document.tex`)

        socket.addEventListener("message", (event) => {

            const message = JSON.parse(event.data);

            switch (message.type) {

                case "user_connected":
                    console.log("New user connected. ID: " + message.id);
                    SetUserID(message.id)

                    socket.send(JSON.stringify({
                        type:     "stateRequest",
                    }));
                    // todo: promise och wait om ingen respons?

                    break;

                case "stateResponse":
                    updateFromCode = true;

                    const jsonString = decodeByteData(message.byteState as string)
                    
                    const loadedDocStr = LoadState(jsonString)

                    console.log("Recieved doc state!")
                    dispatchStateToEditor(loadedDocStr)

                    updateFromCode = false;
                    break;
                    

                case "operation":
                    console.log("Operation recieved")
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
