<script>
    let SavedProjects = 0;

    function handleClick() {
        window.location.href = "/EditorArea";
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
