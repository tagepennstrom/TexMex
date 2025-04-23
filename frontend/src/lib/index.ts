// src/lib/index.ts
import { get } from 'svelte/store';
import { editorView } from '$lib/stores';

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