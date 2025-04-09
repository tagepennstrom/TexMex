<script lang='ts'>
    import Header from './Header.svelte';
    import Viewer from './Viewer.svelte';
    import Editor from './Editor.svelte';


    let pdfUrl = $state("");
    let compileCount = $state(0);

    async function compileLatex(content: string) {
        const serverUrl = `http://${location.hostname}:8080`;
        const res = await fetch(`${serverUrl}/compileDocument`, {
            method: "POST",
            headers: { "Content-Type": "text/plain" },
            body: content
        });

        if (!res.ok) {
            console.error("Failed to compile LaTeX:", await res.text());
            return;
        }

        const data = await res.json();
        pdfUrl = serverUrl + data.pdfUrl;
        compileCount++;
    }
</script>

<Header/>
<Editor {compileLatex} />
<Viewer {pdfUrl} {compileCount}/>
