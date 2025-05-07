<script>
    let { pdfUrl, compileCount } = $props();
</script>

{#key compileCount}
    <div class="viewer-container">
        <iframe
            src={pdfUrl}
            class="pdf-viewer"
            aria-label="PDF Preview"
            title="PDF Viewer"
        ></iframe>
    </div>
{/key}

<style>
    .viewer-container {
        width: 49%; /* Adjust the width of the viewer container */
        height: 700px; /* Set a fixed height for the viewer */
        float: right;
        margin: auto;

    }

    .pdf-viewer {
        width: 100%; /* Make the iframe fill the container */
        height: 100%; /* Make the iframe fill the container */
        border: none; /* Remove the default iframe border */
    }
</style>
