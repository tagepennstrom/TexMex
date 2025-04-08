<script lang=ts>
    export let textarea: HTMLTextAreaElement | null;

    function getCaretPosition(): number {
        if (!textarea) return -1; //null case
        const currentLineIndex = textarea.value.lastIndexOf('\n', textarea.selectionStart - 1) + 1;
        return currentLineIndex;
    }
</script>


<button on:click={getCaretPosition}>Get index</button>

<p>
    Index i textarea: {getCaretPosition()}
</p>

<style>
    p {
        float: left;
    }

    button {
        float: left;
    }
</style>