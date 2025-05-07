import { get } from 'svelte/store';
import { editorView, compileLatexStore } from "$lib/stores";
import { showFilesModal } from '$lib/stores';



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

export function insertNewline(){
    const view = get(editorView);
    if(!view) return;

    const { state, dispatch } = view;
    const cursor = state.selection.main.head;

    dispatch({
        changes: {
            from: cursor,
            to: cursor,
            insert: '\\newline{}'
        }
    });

    dispatch({
        selection: { anchor: cursor + 9, head: cursor + 9 }
    });
}


export function compile() {
    const view = get(editorView);
    const compileLatex = get(compileLatexStore);

    if (view && compileLatex) {
        let content = view.state.doc.toString();

        
        compileLatex(content);
    }
 }
 
export const toggleFilesModal = () => {
    showFilesModal.update(value => !value);
};




