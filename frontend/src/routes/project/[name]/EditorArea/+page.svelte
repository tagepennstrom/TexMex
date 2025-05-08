<script lang='ts'>
    import { page } from '$app/state'
    import Header from '../../../Header.svelte';
    import Viewer from './Viewer.svelte';
    import Editor from './Editor.svelte';
    import Footer from './Footer.svelte';
	import Toolbar from './Toolbar.svelte';


    let pdfUrl = $state("");
    let compileCount = $state(0);

    async function compile() {
        const serverUrl = `http://${location.hostname}:8080`;
        const res = await fetch(`${serverUrl}/projects/${page.params.name}/pdf`);

        if (!res.ok) {
            console.error("Failed to compile LaTeX:", await res.text());
            return;
        }

        compileCount++;

        const pdfBytes = await res.bytes();
        const blob = new Blob([pdfBytes], { type: 'application/pdf' });
        pdfUrl = URL.createObjectURL(blob);
    }
</script>



<div class="page-container">
    <Header/>
    <div class="toolbar">
        <Toolbar {compile} />
    </div>
    <div class="content">
        <Editor />
        <Viewer {pdfUrl} {compileCount}/>
    </div>
    <Footer/>
</div>

<style>
    .page-container {
        display: flex;
        flex-direction: column;
        
    }

    .content {
        flex: 1;
    }
</style>
