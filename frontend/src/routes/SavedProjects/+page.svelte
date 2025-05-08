<script lang='ts'>
    import {onMount} from 'svelte'
    import Header from '../Header.svelte';
    import HiddenToolbar from '../HiddenToolbar.svelte';

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

</script>



<div class="page-container">
    <Header/>
    <div class="hiddenToolbar">
        <HiddenToolbar/>
    </div>
    <div class ="contents">
        {#if projects.length === 0}
            <p>No saved projects</p>
        {:else}
            <ul>
            {#each projects as project}
                <li>
                    <a href='/project/{project.name}/EditorArea'>{project.name}</a>
                </li>
            {/each}
            </ul>
        {/if}
    </div>
</div>

<style>
    .page-container {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }


</style>
