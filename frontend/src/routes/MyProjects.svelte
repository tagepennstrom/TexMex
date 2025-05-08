<script lang="ts">
    import {onMount} from 'svelte'

    type Project = {
        name: string;
        documents: Document[];
    }
    type Document = {
        name: string;
    }

    let projects: Project[] = $state([]);

    onMount(() => {
        const serverUrl = `http://${location.hostname}:8080`;
        fetch(`${serverUrl}/projects`)
            .then(res => res.json())
            .then((res: Project[]) => projects = res);
    });

    function handleClick() {
        window.location.href = "/SavedProjects";
    }

</script>

<div class="MyProjects">
    <div class="projects">
        {#if projects.length === 0}
            <p>No saved projects</p>
        {:else}
            {#each projects as project}
                <li>{project.name}</li>
            {/each}
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
