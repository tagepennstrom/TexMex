<script lang="ts">
    import { insertBold, insertItalic, insertUnderline } from '$lib/index';
 



</script>
<div class="Toolbar">
    <button class="Bold" onclick={insertBold} title ="Bold"><strong>B</strong></button>
    <button class="Italic" onclick={insertItalic} title="Italic"><em>I</em></button>
    <button class="Underline" onclick={insertUnderline} title="Underline"><u>U</u></button>
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


</style>
