<script lang='ts'>
    export let latexContent: string;
    export let onCompile: (content: string) => void; // Update the type to accept content
</script>

<button on:click={() => onCompile(latexContent)}>Compile</button> <!-- Pass the latest content -->
<p></p>
<textarea 
    name="Edit" placeholder="Write LaTex here"
    bind:value={latexContent}
></textarea>

<style>
    textarea {
        width: 49%;
        height: 700px;
        float: left;
        margin: auto;
    }

    button {
        display: block; /* Make the button a block element */
        margin: 10px auto; /* Center the button horizontally */
        padding: 10px 20px;
        background-color: darkorange;
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    button:hover {
        background-color: orange;
    }
</style>