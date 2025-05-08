<script>
    async function handleClick() {
        const serverUrl = `http://${location.hostname}:8080`;
        const projectName = 'project-1';
        const res = await fetch(`${serverUrl}/projects/${projectName}`, {
            method: "POST",
        });
        if (!res.ok) {
            console.error("Unable to create project:", await res.text());
            return;
        }
        window.location.href = `/project/${projectName}/EditorArea`;
    }
</script>

<div class="Footer">
    <button class="button" onclick={handleClick}>Start new project</button>
</div>

<style>
    .Footer {
        display: flex;
        justify-content: center;
        align-items: center; 
    }

    .button{
        align-self: center;
        background: #f0f0f0;
        border: 1px solid black;
        color: black;
        font-size: 1rem;
        padding: 0.5rem 1rem;
        cursor: pointer;
        text-align: center;
        margin-top: 1rem;
    }
</style>
