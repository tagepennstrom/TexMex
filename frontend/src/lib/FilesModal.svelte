<script lang='ts'>
	import { toggleFilesModal } from '$lib';
    import { showFilesModal } from './stores';
    import {fly } from 'svelte/transition';
    export let projectName: String | null = null;

    type AllFiles = {
        name: String
    }

    let files: FileList | null = null;
    let allFiles: AllFiles[];
    let uploadMessage: string = "";
    let uploadMode: boolean = false;

    function toggleUploadMode() {
        uploadMode = !uploadMode
    }

    async function uploadFiles() {
        uploadMessage = "";
        if (!projectName) {
            console.error("Projectname is null")
            return;
        }
  
        if (!files || files.length === 0) {
            console.error("No files selected");
            return;
        }

        const maxSizeMB = 1;
        for (let i = 0; i < files.length; i++) {
            if (files[i].size > maxSizeMB * 1024 * 1024) {
                alert(`File "${files[i].name}" is bigger than ${maxSizeMB} MB and cannot be uploaded.`);
                return;
            }
        }

        const formData = new FormData();
        for (let i = 0; i < files.length; i++) {
            formData.append("file", files[i]);
        }

        try {
            const serverUrl = `http://${location.hostname}:8080`;
            const response = await fetch(`${serverUrl}/projects/uploadFile?projectName=${projectName}`, {
            method: "POST",
            body: formData,
        });

            if (!response.ok) {
                throw new Error("Failed to upload files");
            }
            uploadMessage = "Files uploaded sucessfully"
            console.log("Files uploaded successfully");
        } catch (error) {
            console.error("Error uploading files:", error);
        }
    }

    function removeFile(index: number) {
        if (!files) return;
        const fileArray = Array.from(files);
        fileArray.splice(index, 1);
        const dataTransfer = new DataTransfer();
        fileArray.forEach(file => dataTransfer.items.add(file));
        files = dataTransfer.files;
    }

</script>

{#if $showFilesModal}
<div class="backdrop">
    <div class="modal">
        <button class="close" onclick={toggleFilesModal}>Close</button>
        <button class="close" onclick={toggleUploadMode}>Upload Files</button>
        {#if uploadMode}    
            <input bind:files type="file" id="file" multiple accept=".pdf,.tex,.txt,.bib,image/*"/>
            <label for="file">Upload File</label>
            <button onclick="{uploadFiles}">Upload</button>
            {#if files && !uploadMessage}
                <h2>Selected files:</h2>
                {#each Array.from(files) as file, i}
                    <p>
                        {file.name} ({(file.size / (1024)).toFixed(0)} KB)
                        <button class="remove-file" onclick={() => removeFile(i)} title="Remove file">✖</button>
                    </p>
                    {/each}
                    {/if}
                {#if uploadMessage}
                    <p style="color: green; margin-top: 1em;">{uploadMessage}</p>
                {/if}
        {/if}
    </div>
</div>
{/if}

<style>
    .backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.5);
        z-index: 1000; 
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal {
        background-color: white;
        border-radius: 10px;
        padding: 20px;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
        z-index: 1001;
    }

    Input {
        Display: none;
    }

    Label {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 200px; /* Öka höjden */
        width: 300px; /* Öka bredden */
        border-radius: 20px; /* Justera hörnens rundning */
        border: 2px dashed #999; /* Gör gränsen mer framträdande */
        font-size: 1.5rem; /* Gör texten större */
        text-align: center;
        cursor: pointer; /* Visa en pekare för att indikera klickbarhet */
        background-color: #f9f9f9;
    }

    Label:hover {
        Color: #de0611;
        Border: 1px dashed #de0611;
    }
</style>


