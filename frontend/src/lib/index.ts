// src/lib/index.ts
import { get } from 'svelte/store';
import { editorView, compileLatexStore } from "$lib/stores";



export function insertBold() {
    const view = get(editorView);
    if (!view) return;

    const { state, dispatch } = view;
    const cursor = state.selection.main.head;

    dispatch({
        changes: {
            from: cursor,
            to: cursor,
            insert: '\\textbf{}'
        }
    });

    dispatch({
        selection: { anchor: cursor + 8, head: cursor + 8 }
    });
}

export function insertItalic(){
    const view = get(editorView);
    if(!view) return;

    const { state, dispatch } = view;
    const cursor = state.selection.main.head;

    dispatch({
        changes: {
            from: cursor,
            to: cursor,
            insert: '\\textit{}'
        }
    });

    dispatch({
        selection: { anchor: cursor + 8, head: cursor + 8 }
    });
}

export function insertUnderline(){
    const view = get(editorView);
    if(!view) return;

    const { state, dispatch } = view;
    const cursor = state.selection.main.head;

    dispatch({
        changes: {
            from: cursor,
            to: cursor,
            insert: '\\underline{}'
        }
    });

    dispatch({
        selection: { anchor: cursor + 11, head: cursor + 11 }
    });
}


export function compile() {
    const view = get(editorView);
    const compileLatex = get(compileLatexStore);

    if (view && compileLatex) {
        const content = view.state.doc.toString();
        compileLatex(content);
    }
 }




