<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { insertBold, insertItalic, insertUnderline } from '$lib/index';

    // Event handler for keyboard shortcuts
    function handleKeydown(event: KeyboardEvent) {
        if (event.ctrlKey) { 
            event.preventDefault(); // Prevent default browser action for key combinations

            switch (event.key) {
                case 'b':
                    insertBold();
                    break;
                case 'i':
                    insertItalic();
                    break;
                case 'u':
                    insertUnderline();
                    break;
            }
        }
    }

    // Add and remove the event listener for keydown
    onMount(() => {
        document.addEventListener('keydown', handleKeydown);
    });

    onDestroy(() => {
        document.removeEventListener('keydown', handleKeydown);
    });
</script>

<div class="Toolbar">
    <button class="Bold" onclick={insertBold} title="Bold"><strong>B</strong></button>
    <button class="Italic" onclick={insertItalic} title="Italic"><em>I</em></button>
    <button class="Underline" onclick={insertUnderline} title="Underline"><u>U</u></button>
</div>

<style>
    .Toolbar {
        display: flex;
        gap: 10px;
        padding: 10px;
        background-color: #e0f7fa;
        border-bottom: 1px solid #ccc;
        align-items: center;
    }

    .Bold, .Italic, .Underline {
        padding: 10px 20px;
        background-color: rgb(246, 137, 4);
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    .Bold:hover, .Italic:hover, .Underline:hover {
        background-color: orange;
    }
</style>
