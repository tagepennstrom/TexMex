<script lang='ts'>
    import {basicSetup, EditorView} from "codemirror"
    import {onMount} from 'svelte'
    import {EditorState, Transaction} from "@codemirror/state"
    import {StreamLanguage,} from '@codemirror/language'
    import { stex } from "@codemirror/legacy-modes/mode/stex"
    import { editorView as editorViewStore , compileLatexStore} from "$lib/stores";
    import { autocompletion } from "@codemirror/autocomplete";
    import { myCompletions } from "$lib/completions";


    let { compileLatex } = $props();
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
    
    type Message = {
        document: string
        changes: Change[]
        cursorIndex: number
    }

    type UpdatedDocMessage = {
        document: string
        cursorIndex: number
    }

    async function applyUpdate(document: string, changes: Change[]) {
        updateFromCode = true;
        const updatedDocMessage: UpdatedDocMessage = UpdateDocument(
            document,
            changes,
            editorView.state.selection.main.anchor
        )
        console.log(updatedDocMessage);
        editorView.dispatch({
            changes: {
                from: 0,
                to: editorView.state.doc.length,
                insert: updatedDocMessage.document,
            },
            selection: {
                anchor: updatedDocMessage.cursorIndex,
            },
        });
        updateFromCode = false;
    }

    function sendChangesToCrdt(tr: Transaction): void {
        const cursorIndex = editorView.state.selection.main.anchor;
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
        const updatedDocMessage: UpdatedDocMessage = UpdateDocument(
            document,
            changes,
            editorView.state.selection.main.anchor
        )

        const message: Message = {
            document: document,
            changes: changes,
            cursorIndex: cursorIndex,
        };
        console.log("Sending message:", message);
        socket.send(JSON.stringify(message));

        const serverUrl = `http://${location.hostname}:8080`;
        fetch(`${serverUrl}/saveDocument`, {
            method: "POST",
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
            const message: Message = JSON.parse(event.data);
            console.log(message);
            applyUpdate(editorView.state.doc.toString(), message.changes)
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

                // Spara editorn i store för delning med Toolbar
                editorViewStore.set(editorView);
                compileLatexStore.set(compileLatex);
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
