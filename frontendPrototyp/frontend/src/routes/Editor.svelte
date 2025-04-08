<script lang='ts'>
    import {basicSetup, EditorView} from "codemirror"
    import {onMount} from 'svelte'
    import {EditorState} from "@codemirror/state"
    import {ViewUpdate} from "@codemirror/view"
    import {
        StreamLanguage,
    } from '@codemirror/language'
    import { stex } from "@codemirror/legacy-modes/mode/stex"


    const serverUrl = "http://localhost:8080";
    let socket: WebSocket;

    let latexContent = $state('');
    let { compileLatex } = $props();

    let editor: HTMLElement;
    let editorView: EditorView;

    function onUpdate(update: ViewUpdate) {
        console.log(update);
        const message = {
            document: editorView.state.doc.toString(),
        };
        socket.send(JSON.stringify(message));
    }

    onMount(() => {
        socket = new WebSocket(`${serverUrl}/editDocWebsocket`);

        socket.addEventListener("message", (event) => {
            const res = JSON.parse(event.data);
            console.log(res);
            // TODO: update editorView state
        });

        fetch(`${serverUrl}/document`)
            .then(res => res.text())
            .then(text => {
                // Initialize CodeMirror editor
                latexContent = text;
                editorView = new EditorView({
                    state: EditorState.create({
                        doc: latexContent,
                        extensions: [
                            basicSetup,
                            StreamLanguage.define(stex),
                            EditorView.updateListener.of(onUpdate),
                        ]
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
        display: block;
        margin: 10px auto; /* Center the button horizontally */
        padding: 10px 20px;
        background-color: darkorange;
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    button:hover {
        background-color: orange;
    }
</style>
