<script lang='ts'>
    export let latexContent: string;
    export let onCompile: (content: string) => void; // Update the type to accept content
    import {basicSetup, EditorView} from "codemirror"
    import {onMount} from 'svelte'
    import {EditorState, Compartment} from "@codemirror/state"
    import {
        syntaxHighlighting,
        defaultHighlightStyle,
        StreamLanguage,
    } from '@codemirror/language'
    import { stex } from "@codemirror/legacy-modes/mode/stex"

    let editor: HTMLDivElement | null = null; // Reference to the editor container
    let editorView: EditorView | null = null;


    onMount(() => {
        if (editor) {
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
        }

        return () => {
            editorView?.destroy(); // Clean up editor on component destroy
        };
    });

    function compileContent() {
        if (editorView) {
            const content = editorView.state.doc.toString(); // Get the current content from CodeMirror
            onCompile(content); // Pass the content to the compile function
        }
    }
``
</script>


<button on:click={() => compileContent()}>Compile</button> <!-- Pass the latest content -->
<p></p>

<div class="editor" bind:this={editor}></div>

<!-- <textarea 
    name="Edit" placeholder="Write LaTex here"
    bind:value={latexContent}
    bind:this={textareaRef}
></textarea> -->

<style>
    textarea {
        width: 49%;
        height: 700px;
        float: left;
        margin: auto;
    }
    .editor {
        height: 700px;
        width: 49%;
        float: left;
        margin: auto;
    }

    button {
        display: block; /* Make the button a block element */
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