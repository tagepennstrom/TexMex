<script lang='ts'>
    import Header from './Header.svelte';
    import Viewer from './Viewer.svelte';
    import Editor from './Editor.svelte';
    import SentOutput from './SentOutput.svelte';

    let latexContent = "";
    let pdfUrl: string | null = "http://localhost:8080/pdf";

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

            const res = await fetch("http://localhost:8080/compile", {
                method: "POST",
                headers: { "Content-Type": "text/plain" },
                body: content
            });

            if (res.ok) {
                const data = await res.json();
                if (data.pdfUrl) {
                    // Update the PDF URL with a unique timestamp using Date.now()
                    pdfUrl = `http://localhost:8080${data.pdfUrl}?t=${Date.now()}`;
                    console.log("Updated pdfUrl:", pdfUrl); // Log the updated URL
                } else {
                    console.error("PDF URL missing in response");
                }
            } else {
                console.error("Failed to compile LaTeX:", await res.text());
            }
        } catch (error) {
            console.error("Error during compilation:", error);
        }
    }
</script>

<Header/>

<Editor {latexContent} onCompile={compileLatex} />

<Viewer {pdfUrl}/>