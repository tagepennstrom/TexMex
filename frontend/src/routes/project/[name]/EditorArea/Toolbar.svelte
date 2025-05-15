<script lang="ts">
    import { insertBold, insertItalic, insertUnderline, insertNewline, toggleFilesModal } from '$lib/index';
    import { onMount, onDestroy } from 'svelte';

    let { compile } = $props()

    let handleKeydown: (event: KeyboardEvent) => void;

    onMount(() => {
        handleKeydown = (event: KeyboardEvent) => {

            if (event.ctrlKey) { 
                event.preventDefault();
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
                    case 'l':
                    event.preventDefault(); // Prevents the new window action
                        insertNewline();
                        break;
                }
            }
        };

        document.addEventListener('keydown', handleKeydown);

        onDestroy(() => {
            document.removeEventListener('keydown', handleKeydown);
        });
    });
</script>

<div class="Toolbar">
    
    <div class="button-container">
        <button class="Bold Toolbar-Type" onclick={insertBold} title ="Bold (crtl + b)"><strong>B</strong></button> 
        <button class="Italic Toolbar-Type" onclick={insertItalic} title="Italic (crtl + i)"><em>I</em></button>
        <button class="Underline Toolbar-Type" onclick={insertUnderline} title="Underline (crtl + u)"><u>U</u></button>
        <button class="Newline Toolbar-Type" onclick={insertNewline} title="New line (crtl + l)">\n</button>
    </div>

    <div class="button-container">
        <button class="Toolbar-Type other-features" onclick={toggleFilesModal} title ="Show project files">Files</button>
        <button class="Toolbar-Type other-features" onclick={compile}>Compile</button>
    </div>

</div>

<style>

    .Toolbar{
        display: flex;
        width: 100%;
        background-color: rgb(250, 255, 239);
        border-bottom: 1px solid #ccc;
        box-shadow: 1px 8px 14px -5px rgba(0,0,0,0.15);
        justify-content: space-between;
    }

    .button-container{
        display: flex;
        gap: 10px;
        padding: 10px;
        align-items: center;
        margin: 0 1rem;

    }

    .Toolbar-Type {
        background-color: rgb(123, 197, 53);
        padding: 10px 20px;
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;

    }

    .Toolbar-Type:hover{
        background-color: rgb(99, 160, 43);

    }

    .other-features{
        background-color: rgb(42, 42, 42);
    }

    .text-function{

    }
   


</style>
