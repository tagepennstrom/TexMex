<script lang='ts'>
    import {basicSetup, EditorView} from "codemirror"
    import {onMount} from 'svelte'
    import {EditorState} from "@codemirror/state"
    import {ViewUpdate} from "@codemirror/view"
    import {StreamLanguage,} from '@codemirror/language'
    import { stex } from "@codemirror/legacy-modes/mode/stex"


    let { compileLatex } = $props();
    const serverUrl = "http://localhost:8080";
    let socket: WebSocket;

    let broadcastUpdate = false;

    let editor: HTMLElement;
    let editorView: EditorView;

    
    type Change = {
        from: number;   // Start index
        to: number;     // Slut index
        text: string;   // Tillagd text, tom vid borttagning
    }
    
    type Message = {
        document: string;
        changes: Change
    }

    function onUpdate(update: ViewUpdate) {
        if (!update.docChanged || broadcastUpdate) return;
        
        //Skickar bara det som ändras
        const changes: Change[] = [];
        update.changes.iterChanges((fromA, toA, fromB, toB, inserted) => {
            changes.push({
                from: fromA, 
                to: toA,     
                text: inserted.toString() // Tillagd text, tom vid borttagning
            });
        });
        
        // Skickar hela dokumentet
        const message = {
            document: editorView.state.doc.toString(),
            changes: changes
        };
        console.log("Sending message:", message);
        console.log("Changes:", message.changes);
        socket.send(JSON.stringify(message));
        broadcastUpdate = false;
    }

    
    const fixedHeightEditor = EditorView.theme({
        "&": {height: "700px"},
        ".cm-scroller": {overflow: "auto"}
    })

    const extensions = [
        basicSetup,
        StreamLanguage.define(stex),
        EditorView.updateListener.of(onUpdate),
        fixedHeightEditor
    ]


    onMount(() => {
        socket = new WebSocket(`${serverUrl}/editDocWebsocket`);

        socket.addEventListener("message", (event) => {
            const res = JSON.parse(event.data);
            console.log(event);
            broadcastUpdate = true;
            editorView.dispatch({
                changes: {
                    from: 0,
                    to: editorView.state.doc.length,
                    insert: res.document,
                }
            });
        });

        fetch(`${serverUrl}/document`)
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
            });
    });


    function compileContent() {
        const content = editorView.state.doc.toString(); // Get the current content from CodeMirror
        compileLatex(content);
    }
</script>


<button onclick={() => compileContent()}>Compile</button>
<div id="editor" bind:this={editor}></div>

<style>
    #editor {
        height: 700px;
        width: 49%;
        float: left;
        margin: auto;
    }

    button {
        padding: 10px 20px;
        background-color: darkorange;
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        align-items:flex-start;
    }

    button:hover {
        background-color: orange;
    }
</style>
