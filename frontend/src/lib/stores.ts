// src/lib/stores.ts
import { writable } from 'svelte/store';
import type { EditorView } from '@codemirror/view';

export const editorView = writable<EditorView | null>(null);
export const compileLatexStore = writable<(content: string) => void>();
export const showFilesModal = writable(false);
