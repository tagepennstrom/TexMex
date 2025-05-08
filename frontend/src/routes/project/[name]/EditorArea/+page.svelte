    <script lang='ts'>
        import { page } from '$app/state'
        import Header from '../../../Header.svelte';
        import Viewer from './Viewer.svelte';
        import Editor from './Editor.svelte';
        import Footer from './Footer.svelte';
        import Toolbar from './Toolbar.svelte';


        let compileError = $state(0);    
        let pdfUrl = $state("");
        let compileCount = $state(0);
        let errorMessage = $state("");

    function extractErrorsUsingRegex(logText: string): string[] {
        
        const regex = /(^!.*(?:\n(?!\s*$).*)*)/gm;
        const matches = [];
        let match;
        let len = 0;
        
        while ((match = regex.exec(logText)) !== null) {
        let errorMessage = match[0].trim();

        matches.push(errorMessage);
    }
    
        //ta bort de två sista elementen av arrayn då de inte är användbara
        if(matches[len - 1] == "!  ==> Fatal error occurred, no output PDF file produced!"){
            compileError = 1;
            delete matches[len-2];
            delete matches[len-1];
        }else{
            compileError = 2;
        }
        return matches;
    }

    async function compile() {
        compileError = 0;
        errorMessage = ""
        const serverUrl = `http://${location.hostname}:8080`;
        const res = await fetch(`${serverUrl}/projects/${page.params.name}/pdf`);

        

        if (!res.ok) {
            //fetch projects logFile
            const logText = await fetch(`${serverUrl}/projects/${page.params.name}/documents/document.log`)
                    .then(response => response.text())
            console.log(logText);
                    
            console.error("Failed to compile LaTeX:", await res.text());
            const errorArray = extractErrorsUsingRegex(logText);

            let counter = 1;
            errorArray.forEach(message => {
                errorMessage += `Error ${counter}: `;
                errorMessage += message;
                errorMessage += "\n\n";
                counter += 1;

            });
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
            {#if compileError === 1 || compileError === 2}
                <div class="error-console">
                    <strong>
                        {compileError === 1 
                            ? 'Compiling error, but could compile:' 
                            : 'Compiling error, could not compile:'}
                    </strong>
                    <pre>{errorMessage}</pre>
                </div>
            {/if}
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

        .error-console {
            background: #ffefef;
            color: #a00;
            padding: 1em;
            margin: 1em 0;
            border: 1px solid #a00;
            border-radius: 6px;
            font-family: monospace;
        }
    </style>
