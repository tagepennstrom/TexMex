<script lang='ts'>
    import Header from './Header.svelte';
    import Viewer from './Viewer.svelte';
    import Editor from './Editor.svelte';

    let latexContent = $state("");
    const pdfUrl = "http://localhost:8080/pdf";
    let compileCount = $state(0);

    // Load saved content from localStorage only in the browser
    if (typeof window !== "undefined") {
        latexContent = localStorage.getItem("latexContent") || "";
    }

    async function compileLatex(content: string) {
        try {
            // Save the current content to localStorage only in the browser
            if (typeof window !== "undefined") {
                localStorage.setItem("latexContent", content);
            }

            const res = await fetch("http://localhost:8080/compileDocument", {
                method: "POST",
                headers: { "Content-Type": "text/plain" },
                body: content
            });
            compileCount++;

        } catch (error) {
            console.error("Error during compilation:", error);
        }
    }
</script>

<Header/>

<Editor {latexContent} onCompile={compileLatex} />

<Viewer {pdfUrl} {compileCount}/>
