<script lang='ts'>
    import Header from './Header.svelte';
    import Viewer from './Viewer.svelte';
    import Editor from './Editor.svelte';
    import Footer from './Footer.svelte';


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



<div class="page-container">
    <Header/>
    <div class="content">
        <div class="toolbar">
            <Toolbar/>
        </div>
        <button class="compile-button" onclick={compileLatex(content)}>Compile</button>
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

    .compile-button {
        float: top;
    }
    

</style>

