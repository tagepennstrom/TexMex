<script lang='ts'>
	import { toggleFilesModal } from '$lib';
    import { showFilesModal } from './stores';
    import {fly } from 'svelte/transition';

    type AllFiles = {
        name: String
    }

    let files: FileList | null = null;
    export let projectName: String | null = null;
    let allFiles: AllFiles[];

    async function uploadFiles() {
        if (!files || files.length === 0) {
            console.error("No files selected");
            return;
        }

        const formData = new FormData();
        for (let i = 0; i < files.length; i++) {
            formData.append("file", files[i]);
        }

        try {
            const serverUrl = `http://${location.hostname}:8080`;
            const response = await fetch(`${serverUrl}/projects/${projectName}/uploadFile`, {
            method: "POST",
            body: formData,
        });

            if (!response.ok) {
                throw new Error("Failed to upload files");
            }

            console.log("Files uploaded successfully");
        } catch (error) {
            console.error("Error uploading files:", error);
        }
    }

</script>

{#if $showFilesModal}
<div class="backdrop">
    <div class="modal">
        <button class="close" onclick={toggleFilesModal}>Close</button>
        <input bind:files type="file" id="file"/>
        <label for="file">Upload File</label>
        <button onclick="{uploadFiles}">Upload</button>
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


