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
    
        <div class="title-container .outfit-300">
            <h1 class="cal-sans-regular">All projects</h1>
        </div>


        <div class="meta-data-header-container">
            <div class="meta-data-head">Name</div>
            <div class="meta-data-head">Last modified</div>
        </div>

    <hr id='projects-separator'>


        {#if projects.length === 0}
            <p>No saved projects</p>
        {:else}
            <ul>
            {#each projects as project}
                <a href='/project/{project.name}/EditorArea' class="project-container">
                    <svg xmlns="http://www.w3.org/2000/svg" class="docIcon" viewBox="0 0 24 24" fill="none">
                        <path fill-rule="evenodd" clip-rule="evenodd" d="M9.29289 1.29289C9.48043 1.10536 9.73478 1 10 1H18C19.6569 1 21 2.34315 21 4V20C21 21.6569 19.6569 23 18 23H6C4.34315 23 3 21.6569 3 20V8C3 7.73478 3.10536 7.48043 3.29289 7.29289L9.29289 1.29289ZM18 3H11V8C11 8.55228 10.5523 9 10 9H5V20C5 20.5523 5.44772 21 6 21H18C18.5523 21 19 20.5523 19 20V4C19 3.44772 18.5523 3 18 3ZM6.41421 7H9V4.41421L6.41421 7ZM7 13C7 12.4477 7.44772 12 8 12H16C16.5523 12 17 12.4477 17 13C17 13.5523 16.5523 14 16 14H8C7.44772 14 7 13.5523 7 13ZM7 17C7 16.4477 7.44772 16 8 16H16C16.5523 16 17 16.4477 17 17C17 17.5523 16.5523 18 16 18H8C7.44772 18 7 17.5523 7 17Z" fill="#000000"/>
                    </svg>

                    <div class="project-textcontent">
                        <div class="outfit-500">{project.name}</div>
                        <div class="outfit-300">Apr 20, 2025</div>
                    </div>


                </a>

            {/each}
            </ul>
        {/if}
    </div>
</div>

<style>

    @import url('https://fonts.googleapis.com/css2?family=Cal+Sans&display=swap');
    @import url('https://fonts.googleapis.com/css2?family=Outfit:wght@100..900&display=swap');

    .cal-sans-regular {
        font-family: 'Cal Sans', sans-serif;
        font-weight: 100;
        font-style: normal;
    }

    .outfit-300 {
        font-family: "Outfit", sans-serif;
        font-optical-sizing: auto;
        font-weight: 300;
        font-style: normal;
    }

    .outfit-500 {
        font-family: "Outfit", sans-serif;
        font-optical-sizing: auto;
        font-weight: 500;
        font-style: normal;
    }

    .title-container{
        display: flex;
        flex-direction: row;
        align-items: center;
        width: 90%;
        justify-content:space-between;
        margin-top: 1rem;
        margin-bottom: 3rem;
        white-space: nowrap;
        margin-left: 36px;

    }

    .add-button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        background-color: #34a93c;
        color: white;
        border-radius: 999px;
        padding: 0;
        width: 60px;
        height: 60px;
        overflow: hidden;
        transition: width 0.3s ease, padding 0.3s ease;
        cursor: pointer;
        white-space: nowrap;
        box-shadow: rgba(0, 0, 0, .2) 0 2px 4px 0;
        border: none;

    }

    .add-button .label {
        opacity: 0;
        margin-left: 0;
        width: 0;

        transition: opacity 0.3s ease, margin-left 0.3s ease;
        }

    .add-button:hover {
        width: 150px;
        padding: 0 16px;
        }

    .add-button:hover .label {
        opacity: 1;
        width: fit-content;
        margin-left: 8px;
        }


    .meta-data-header-container{
        margin-left: 36px;
        padding-left: 4rem;
        width: 35rem;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        font-family: "Outfit", sans-serif;
        font-weight: 200;
        font-size: 12px;
    }

    .meta-data-head{
        cursor: pointer;
        background-color: #e6e6e6;
        padding: 0.3rem 1rem;
        border-radius: 1rem;
    }

    #projects-separator {
        border: none;
        border-top: 2px solid #393939;
        margin: 10px 0 0 0;
        width: 43rem;
        margin-left: 36px;
    }

    .icon {
        font-size: 20px;
        font-weight: bold;
        }
    
    .project-container{
        display: flex;
        flex-direction: row;
        background-color: rgb(255, 255, 255);
        margin-bottom: 1rem;
        align-items: center;
        padding: 1rem;
        border-radius: 4px;
        width: 40rem;
        box-shadow: rgba(0, 0, 0, .2) 0 2px 4px 0;
        color: inherit;
        text-decoration: none; 

    }

    .project-container:hover{
        cursor: pointer;
        background-color: #e0e0e0;
    }

    .project-textcontent{
        width: 87%;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }

    .docIcon{
        width: 30px;
        height: 30px; 
        margin-right: 1rem;
    }


    .page-container {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }


</style>
