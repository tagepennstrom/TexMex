    <script lang='ts'>
        import { page } from '$app/state'
        import Header from '../../../Header.svelte';
        import Viewer from './Viewer.svelte';
        import Editor from './Editor.svelte';
        import Footer from './Footer.svelte';
        import Toolbar from './Toolbar.svelte';
	import { showFilesModal } from '$lib/stores';
	import FilesModal from '$lib/FilesModal.svelte';
	import { onMount } from 'svelte';


    let compileError = $state(0);    
    let pdfUrl = $state("");
    let compileCount = $state(0);
    let errorMessage : string[] = $state([]);
    let currentErrorIndex = $state(0);
    let projectName = $state<String | null>(null);

    onMount(() => {
        projectName = page.params.name;
        console.log("Project name is: ", projectName);
    })

    function messageExtracting(logText: string): string[] {
        
        const regex = /(^!.*(?:\n(?!\s*$).*)*)|(^Runaway argument\?.*(?:\n(?!\s*$).*)*)/gmi;
        let matches = [];
        let match;
        let len = 0;
        currentErrorIndex = 0;
            
        while ((match = regex.exec(logText)) !== null) {
            let errorMessage = match[0].trim();

            len = matches.push(errorMessage);
        }
        
        console.log(matches[matches.length -1])
        //ta bort det sista elementen av arrayn då det inte är användbart
        if (matches[len -1] == "!  ==> Fatal error occurred, no output PDF file produced!") {
            compileError = 1;
            matches.pop();  // Remove the last element
            
        } else {
            compileError = 2;
        }

        matches = messageCleanUp(matches);
        return matches;
    }

    function messageCleanUp(errorMessages: string[]): string[] {
        let newMessages = errorMessages;
       
        //Manipulera varje string individuellt
        for (let i = 0; i < newMessages.length; i++){
            let temp = newMessages[i];

            //syntax fel
            if(temp.includes("! Undefined control sequence.")){
                let brackets = false;
                temp = temp.replace("! Undefined control sequence.", "Incorrect LaTeX syntax");
                //Kollar om orginal meddelandet innehåller {}
                if(temp.includes("{}")){
                    brackets = true 
                }

                //Byter ut 1.4 t.e.x till bara "found at line 4"
                temp = temp.replace(/^l\.(\d+)/m, "found at line $1:");

                const match = temp.match(/(found at line \d+:.*?\\\w+)/s);
                //Tar bort allt efter \(funktion)
                if (match) {
                    temp = match[1]; 
                }
                //Lägger tillbaka {}
                temp += brackets === true ? "{}" : "";            
            }

            //fel användning av end{}
            if(temp.includes("Runaway argument?")){
                temp = temp.replace(/Runaway argument\?\n\{(\w+)/m, 
                "You forgot to close \\end{$1\nIt should be: \\end{$1}");

                const lines = temp.split("\n");
                temp = lines.slice(0, 2).join("\n");
                //temp = "Somewhere you forgot to close your end function."
            }

            //Inte end{document}
            if(temp.includes("! Emergency stop.")){
                temp = "Document was not ended correctly, either forgotten or misspelled \nshould be \\end{document}"
            }
            newMessages[i] = temp;
        }
        return newMessages;
    }

    async function compile() {
        compileError = 0;
        errorMessage = []
        const serverUrl = `http://${location.hostname}:8080`;
        const res = await fetch(`${serverUrl}/projects/${projectName}/pdf`);

        

        if (!res.ok) {
            //fetch projects logFile
            const logText = await fetch(`${serverUrl}/projects/${projectName}/documents/document.log`)
                    .then(response => response.text())
                    
            console.error("Failed to compile LaTeX:", await res.text());
            errorMessage = messageExtracting(logText);

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
        {#if showFilesModal}
            <div>
                <FilesModal {projectName}/>
            </div>
        {/if}
        <div class="toolbar">
            <Toolbar {compile} />
               {#if errorMessage.length > 0}
            <div class="error-console">
                <strong>
                    Compiling error {errorMessage.length > 1 ? `(Error ${currentErrorIndex + 1} of ${errorMessage.length})` : ''}:
                </strong>
                <pre>{errorMessage[currentErrorIndex]}</pre>

                {#if errorMessage.length > 1}
                    <div class="error-navigation">
                        <!-- svelte-ignore event_directive_deprecated -->
                        <button on:click={currentErrorIndex = (currentErrorIndex - 1 + errorMessage.length) % errorMessage.length}>
                            &lt; <u>Prev</u>
                        </button>
                        <!-- svelte-ignore event_directive_deprecated -->
                        <button on:click={() => currentErrorIndex = (currentErrorIndex + 1) % errorMessage.length}>
                            <u>Next</u> &gt;
                        </button>
                    </div>
                {/if}
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


        .error-navigation button {
            border-radius: 5px; /* You can adjust this number */
            padding: 0.4em 0.8em;
            border: 1px solid #a00;
            background-color: rgba(170, 0, 0, 0);
            color: #a00;
            cursor: pointer;
            margin: 0 0.5em;
            font-family: monospace;
        }

    </style>
