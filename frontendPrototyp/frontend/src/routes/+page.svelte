<script lang='ts'>
    import Header from './Header.svelte';
    import Viewer from './Viewer.svelte';
    import Editor from './Editor.svelte';


    let pdfUrl = $state("");
    let compileCount = $state(0);

    async function compileLatex(content: string) {
        const res = await fetch("http://localhost:8080/compileDocument", {
            method: "POST",
            headers: { "Content-Type": "text/plain" },
            body: content
        });

        if (!res.ok) {
            console.error("Failed to compile LaTeX:", await res.text());
            return;
        }

        const data = await res.json();
        pdfUrl = `http://localhost:8080${data.pdfUrl}`;
        compileCount++;
    }
</script>

<Header/>
<Editor {compileLatex} />
<Viewer {pdfUrl} {compileCount}/>
