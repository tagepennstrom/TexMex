<script lang='ts'>
    import {basicSetup, EditorView} from "codemirror"
    import {onMount} from 'svelte'
    import {EditorState, Transaction} from "@codemirror/state"
    import {StreamLanguage,} from '@codemirror/language'
    import { stex } from "@codemirror/legacy-modes/mode/stex"
    import { editorView as editorViewStore , compileLatexStore} from "$lib/stores";




    let { compileLatex } = $props();
    let socket: WebSocket;

    let broadcastUpdate = $state(false);

    let editor: HTMLElement;
    let editorView: EditorView;

    
    type Change = {
        from: number;   // Start index
        to: number;     // Slut index
        text: string;   // Tillagd text, tom vid borttagning
    }
    
    type Message = {
        changes: Change[]
    }

    function sendChangesToCrdt(tr: Transaction): void {
        const changes: Change[] = [];
            tr.changes.iterChanges((fromA, toA, fromB, toB, inserted) => {
                changes.push({
                    from: fromA, 
                    to: toA,     
                    text: inserted.toString() // Tillagd text, tom vid borttagning
                });
            });

            const message: Message = {
                changes: changes
            };
            console.log("Sending message:", message);
            socket.send(JSON.stringify(message));
    }

    const BlockLocalChanges = EditorState.transactionFilter.of(tr => {
        if (tr.docChanged && !broadcastUpdate) {
            sendChangesToCrdt(tr);
            return []; // Blocka ändringar lokalt genom att skicka tillbaka inga
        }
        return tr; // Inga ändringar var gjorda men vi måste returna.
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
        EditorView.lineWrapping
    ]


    onMount(() => {
        const serverUrl = `http://${location.hostname}:8080`;
        socket = new WebSocket(`${serverUrl}/editDocWebsocket`);

        socket.addEventListener("message", (event) => {
            const res: Message = JSON.parse(event.data);
            console.log(res);
            broadcastUpdate = true;
            res.changes.forEach((change) => {
                editorView.dispatch({
                    changes: {
                        from: change.from,
                        to: change.to,
                        insert: change.text,
                    },
                    selection: {
                        anchor: change.to + 1,
                    },
                });
            });
            broadcastUpdate = false;
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
