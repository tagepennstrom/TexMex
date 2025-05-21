<script lang='ts'>
	import { toggleFilesModal } from '$lib';
    import { showFilesModal } from './stores';
    export let projectName: String | null = null;

    type AllFiles = {
        name: String
    }

    let files: FileList | null = null;
    let allFiles: AllFiles[];
    let uploadMessage: string = "";
    let uploadMode: boolean = false;
    let projectFiles: string[] = [];
    const serverUrl = `http://${location.hostname}:8080`;
    function toggleUploadMode() {
        uploadMode = !uploadMode
    }

    async function delFile(fileName: string) {        
        const response = await fetch(`${serverUrl}/projects/delFile?projectName=${projectName}&fileName=${encodeURIComponent(fileName)}`);
        if (response.ok) {
            console.log("File succesfully deleted");
            fetchProjectFiles();
        }
    }

    async function fetchProjectFiles() {
        if (!projectName) return;
        const response = await fetch(`${serverUrl}/projects/getfiles?projectName=${projectName}`);
        if (response.ok) {
            projectFiles = await response.json();
            console.log("Collected all files");
        } else {
            projectFiles = [];
            console.log("No files in this project");
        }
    }

    $: if ($showFilesModal) {
        fetchProjectFiles();
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
            const response = await fetch(`${serverUrl}/projects/uploadFile?projectName=${projectName}`, {
            method: "POST",
            body: formData,
        });

            if (!response.ok) {
                throw new Error("Failed to upload files");
            }
            uploadMessage = "Files uploaded sucessfully"
            fetchProjectFiles();
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
        <div class="modal-header">
            <h2>Project Files</h2>
            <button class="icon-btn close-btn" onclick={toggleFilesModal} title="Close">✖</button>
        </div>
        <div class="modal-content">
            <button class="toolbar-btn" onclick={toggleUploadMode}>
                {uploadMode ? "Cancel" : "Upload files"}
            </button>

            <h3>Files in project:</h3>
            {#if projectFiles.length === 0}
                <p class="empty-text">No files in this project.</p>
            {:else}
                <ul class="file-list">
                    {#each projectFiles as name}
                        <li class="file-list-item">
                            <span class="file-name">{name}</span>
                            {#if name !== "document.tex"}
                                <button class="icon-btn remove-btn" onclick={() => delFile(name)} title="Delete file">✖</button>
                            {/if}
                        </li>
                    {/each}
                </ul>
            {/if}

            {#if uploadMode}    
                <div class="upload-section">
                    <input bind:files type="file" id="file" multiple accept=".pdf,.tex,.txt,.bib,image/*" class="file-input"/>
                    <label for="file" class="file-label">Choose files to upload</label>
                    <button class="toolbar-btn" onclick={uploadFiles}>Upload</button>
                </div>
                {#if files && !uploadMessage}
                    <div class="selected-files">
                        <h4>Selected files:</h4>
                        {#each Array.from(files) as file, i}
                            <div class="selected-file">
                                <span>{file.name} ({(file.size / 1024).toFixed(0)} KB)</span>
                                <button class="icon-btn remove-btn" onclick={() => removeFile(i)} title="Remove file">✖</button>
                            </div>
                        {/each}
                    </div>
                {/if}
                {#if uploadMessage}
                    <p class="success-message">{uploadMessage}</p>
                {/if}
            {/if}
        </div>
    </div>
</div>
{/if}

<style>
.backdrop {
    position: fixed;
    top: 0; left: 0; width: 100vw; height: 100vh;
    background: rgba(0,0,0,0.35);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.modal {
    background: #fff;
    border-radius: 16px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.18);
    width: 420px;
    max-width: 95vw;
    padding: 0;
    overflow: hidden;
    animation: fadeIn 0.2s;
}

@keyframes fadeIn {
    from { opacity: 0; transform: scale(0.98);}
    to { opacity: 1; transform: scale(1);}
}

.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: rgb(250, 255, 239);
    padding: 1rem 1.5rem;
    border-bottom: 1px solid #e0e0e0;
}

.modal-header h2 {
    margin: 0;
    font-size: 1.3rem;
    color: #2a2a2a;
}

.modal-content {
    padding: 1.5rem;
}

.toolbar-btn {
    background-color: rgb(123, 197, 53);
    color: white;
    border: none;
    border-radius: 5px;
    padding: 0.6em 1.2em;
    font-size: 1rem;
    cursor: pointer;
    margin-bottom: 1em;
    transition: background 0.2s;
}
.toolbar-btn:hover {
    background-color: rgb(99, 160, 43);
}

.icon-btn {
    background: none;
    border: none;
    color: #de0611;
    font-size: 1.2em;
    cursor: pointer;
    padding: 0.2em 0.5em;
    border-radius: 3px;
    transition: background 0.2s, color 0.2s;
}
.icon-btn:hover {
    background: #f7eaea;
    color: #a00000;
}

.close-btn {
    font-size: 1.4em;
    color: #888;
    margin-left: 0.5em;
}
.close-btn:hover {
    color: #de0611;
    background: #f7eaea;
}

.file-list {
    list-style: none;
    padding: 0;
    margin: 0 0 1em 0;
}
.file-list-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.4em 0.2em;
    border-bottom: 1px solid #f0f0f0;
    font-size: 1.05em;
}
.file-name {
    flex: 1;
    color: #222;
    word-break: break-all;
}
.remove-btn {
    margin-left: 0.5em;
}

.empty-text {
    color: #888;
    font-style: italic;
    margin: 1em 0;
}

.upload-section {
    display: flex;
    flex-direction: column;
    gap: 0.7em;
    margin-bottom: 1em;
}

.file-input {
    display: none;
}
.file-label {
    display: inline-block;
    padding: 1em 1.5em;
    border: 2px dashed #999;
    border-radius: 10px;
    background: #f9f9f9;
    color: #555;
    font-size: 1.1em;
    text-align: center;
    cursor: pointer;
    margin-bottom: 0.5em;
    transition: border 0.2s, color 0.2s;
}
.file-label:hover {
    color: #de0611;
    border-color: #de0611;
}

.selected-files {
    margin-top: 1em;
}
.selected-file {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.2em 0;
    font-size: 0.98em;
}
.success-message {
    color: rgb(123, 197, 53);
    margin-top: 1em;
    font-weight: bold;
}
</style>


