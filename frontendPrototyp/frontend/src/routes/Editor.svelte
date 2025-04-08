<script lang='ts'>
    import {basicSetup, EditorView} from "codemirror"
    import {onMount} from 'svelte'
    import {EditorState} from "@codemirror/state"
    import {
        StreamLanguage,
    } from '@codemirror/language'
    import { stex } from "@codemirror/legacy-modes/mode/stex"

    let latexContent = $state("");
    let { compileLatex } = $props();

    let editor: HTMLDivElement;
    let editorView: EditorView;

    onMount(() => {
        // Load saved content from localStorage only in the browser
        latexContent = localStorage.getItem("latexContent") || "";

        // Initialize CodeMirror editor
        editorView = new EditorView({
            state: EditorState.create({
                doc: latexContent,
                extensions: [
                    basicSetup,
                    StreamLanguage.define(stex)
                ]
            }),
            parent: editor
        });
    });

    function compileContent() {
        const content = editorView.state.doc.toString(); // Get the current content from CodeMirror
        compileLatex(content);
    }
</script>


<button onclick={() => compileContent()}>Compile</button>
<div class="editor" bind:this={editor}></div>

<style>
    .editor {
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
