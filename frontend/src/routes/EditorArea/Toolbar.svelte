<script lang="ts">
    import { insertBold, insertItalic, insertUnderline, compile, insertNewline } from '$lib/index';
    import { onMount, onDestroy } from 'svelte';

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
    <button class="Bold" onclick={insertBold} title ="Bold (crtl + b)"><strong>B</strong></button>
    <button class="Italic" onclick={insertItalic} title="Italic (crtl + i)"><em>I</em></button>
    <button class="Underline" onclick={insertUnderline} title="Underline (crtl + u)"><u>U</u></button>
    <button class="Newline" onclick={insertNewline} title="New line (crtl + l)">\n</button>
    <button class="Compile" onclick={compile}>Compile</button>

</div>

<style>

    .Toolbar{
    display: flex;
    gap: 10px;
    padding: 10px;
    background-color: #e0f7fa;
    border-bottom: 1px solid #ccc;
    align-items: center;
    }

    .Bold {
        padding: 10px 20px;
        background-color: rgb(246, 137, 4);
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    .Bold:hover{
        background-color: orange;
    }

    .Italic {
        padding: 10px 20px;
        background-color: rgb(246, 137, 4);
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    .Italic:hover{
        background-color: orange;
    }

    .Underline {
        padding: 10px 20px;
        background-color: rgb(246, 137, 4);
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    .Underline:hover{
        background-color: orange;
    }

    .Compile {
        padding: 10px 20px;
        background-color: rgb(246, 137, 4);
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }
    .Compile:hover{
        background-color: orange;
    }

    .Newline {
        padding: 10px 20px;
        background-color: rgb(246, 137, 4);
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }
    .Newline:hover{
        background-color: orange;
    }


</style>
