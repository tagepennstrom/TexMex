<script lang="ts">
    import { basicSetup, EditorView } from "codemirror";
    import { onMount } from "svelte";
    import { EditorState } from "@codemirror/state";
    import { ViewUpdate } from "@codemirror/view";
    import { StreamLanguage } from "@codemirror/language";
    import { stex } from "@codemirror/legacy-modes/mode/stex";
    import { get } from "svelte/store";
    import { compileLatexStore, editorView as editorViewStore } from "$lib/stores";

    let { compileLatex } = $props();
    const serverUrl = "http://localhost:8080";

    let socket: WebSocket;
    let broadcastUpdate = false;
    let editor: HTMLElement;
    let editorView: EditorView;

    function onUpdate(update: ViewUpdate) {
        if (!update.docChanged || broadcastUpdate) return;

        const message = {
            document: editorView.state.doc.toString(),
        };
        socket.send(JSON.stringify(message));
        broadcastUpdate = false;
    }

    const fixedHeightEditor = EditorView.theme({
        "&": { height: "700px" },
        ".cm-scroller": { overflow: "auto" }
    });

    const extensions = [
        basicSetup,
        StreamLanguage.define(stex),
        EditorView.updateListener.of(onUpdate),
        fixedHeightEditor
    ];

    onMount(() => {
        socket = new WebSocket(`${serverUrl}/editDocWebsocket`);

        socket.addEventListener("message", (event) => {
            const res = JSON.parse(event.data);
            broadcastUpdate = true;

            const view = get(editorViewStore);
            if (view) {
                view.dispatch({
                    changes: {
                        from: 0,
                        to: view.state.doc.length,
                        insert: res.document
                    }
                });
            }
        });

        fetch(`${serverUrl}/document`)
            .then(res => res.text())
            .then(text => {
                editorView = new EditorView({
                    state: EditorState.create({
                        doc: text,
                        extensions
                    }),
                    parent: editor
                });

                // Spara editorn i store för delning med Toolbar
                editorViewStore.set(editorView);
                compileLatexStore.set(compileLatex);
            });
    });

</script>

<!-- UI -->

<div id="editor" bind:this={editor}></div>

<style>
    #editor {
        height: 700px;
        width: 49%;
        float: left;
        margin: auto;
    }

</style>
