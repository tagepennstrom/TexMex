<script>
    let SavedProjects = 0;

    function handleClick() {
        window.location.href = "/SavedProjects";
    }
</script>

<div class="MyProjects">
    <div class="projects">
        {#if SavedProjects > 0}
            <!-- Lägg till sätt att visa sparade projekt -->
            <p>Gör nått som visar alla project</p>
        {:else}
            <p>No saved projects</p>
        {/if}
    </div>
    <div class="footer">
        <button class="button" onclick={handleClick}>Show all my projects</button>
    </div>
</div>

<style>
    .MyProjects {
        height: 800px;
        width: 49%;
        float: left;
        margin: auto;
        border: 1px solid black;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        background: white;
    }

    .projects {
        flex-grow: 1;
        width: 100%;
        color: black;
        padding: 1rem;
    }

    .footer{
        text-align: center;
        margin-bottom: 10px;
    }

    .button{
        background: #f0f0f0;
        border: 1px solid black;
        color: black;
        font-size: 1rem;
        padding: 0.5rem 1rem;
        cursor: pointer;
        text-align: center;
    }
</style>
