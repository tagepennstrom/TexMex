<script lang='ts'>
    import Header from './Header.svelte';
    import Viewer from './Viewer.svelte';
    import Editor from './Editor.svelte';
    import Footer from './Footer.svelte';
	import Toolbar from './Toolbar.svelte';
    import FilesModal from '$lib/FilesModal.svelte';


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



<div class="page-container">
    <Header/>
    <div class="content">
        <div class="toolbar">
            <Toolbar/>
            <FilesModal/>
        </div>
        <Editor {compileLatex} />
        <Viewer {pdfUrl} {compileCount}/>
    </div>
    <Footer/>
</div>

<style>
    .page-container {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }

    .content {
        flex: 1;
    }
</style>

