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

    let showCreateProjectPopup = $state(false);
    let newProjectName = $state('');
    let createProjectPopupError = $state('');

    onMount(() => {
        const serverUrl = `http://${location.hostname}:8080`;
        fetch(`${serverUrl}/projects`)
            .then(res => res.json())
            .then((res: Project[]) => projects = res);
    });

    function goToAllDocumentsPage() {
        window.location.href = "/SavedProjects";
    }

    async function createNewProject() {
        createProjectPopupError = '';
        const serverUrl = `http://${location.hostname}:8080`;
        const res = await fetch(`${serverUrl}/projects/${newProjectName}`, {
            method: "POST",
        });
        if (!res.ok) {
            const errMsg = await res.text();
            createProjectPopupError = errMsg;
            return;
        } 
        const newProjectUrl = `/project/${newProjectName}/EditorArea`; 
        newProjectName = '';
        showCreateProjectPopup = false;
        window.location.href = newProjectUrl;
    }

    function closeCreateProjectPopup() {
        showCreateProjectPopup = false;
        newProjectName = '';
        createProjectPopupError = '';
    }

    function focusElem(elem: HTMLElement) {
        elem.focus(); 
    }
</script>

<div class="MyProjects">
    <div class="title-container">
        <h1 class="cal-sans-regular">Latest projects</h1>
        <button class="add-button outfit-500" onclick={() => showCreateProjectPopup = true}>
            <span class="icon">+</span>
            <span class="label">New project</span>
        </button>

    </div>

    {#if showCreateProjectPopup}
        <div id='create-new-project-popup-overlay'>
            <div id='create-new-project-popup' class='outfit-500'>
                <div id='header'>
                    <h1>New Project</h1>
                    <svg id='closeBtn' onclick={closeCreateProjectPopup} class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                      <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m15 9-6 6m0-6 6 6m6-3a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
                    </svg>
                </div>
                <hr>
                <input type='text' placeholder='Project Name' bind:value={newProjectName} required use:focusElem>
                {#if createProjectPopupError !== ''}
                    <div id='project-name-error-msg'>
                        <span>{createProjectPopupError}</span>
                    </div>
                {/if}
                <hr>
                <button onclick={createNewProject} disabled={newProjectName === ''} id='submitBtn' class='outfit-300'>Create</button>
            </div>
        </div>
    {/if}

    <hr id='projects-separator'>

    <div class="projects">
        {#if projects.length === 0}
            <p>No saved projects</p>
        {:else}
            <ul>
            {#each projects as project}

                <a href='/project/{project.name}/EditorArea' class="project-container">
                    <svg xmlns="http://www.w3.org/2000/svg" class="docIcon" viewBox="0 0 24 24" fill="none">
                        <path fill-rule="evenodd" clip-rule="evenodd" d="M9.29289 1.29289C9.48043 1.10536 9.73478 1 10 1H18C19.6569 1 21 2.34315 21 4V20C21 21.6569 19.6569 23 18 23H6C4.34315 23 3 21.6569 3 20V8C3 7.73478 3.10536 7.48043 3.29289 7.29289L9.29289 1.29289ZM18 3H11V8C11 8.55228 10.5523 9 10 9H5V20C5 20.5523 5.44772 21 6 21H18C18.5523 21 19 20.5523 19 20V4C19 3.44772 18.5523 3 18 3ZM6.41421 7H9V4.41421L6.41421 7ZM7 13C7 12.4477 7.44772 12 8 12H16C16.5523 12 17 12.4477 17 13C17 13.5523 16.5523 14 16 14H8C7.44772 14 7 13.5523 7 13ZM7 17C7 16.4477 7.44772 16 8 16H16C16.5523 16 17 16.4477 17 17C17 17.5523 16.5523 18 16 18H8C7.44772 18 7 17.5523 7 17Z" fill="#000000"/>
                        </svg>


                    <div class="outfit-500">{project.name}</div>
                </a>

               <!--  <li>
                    <a href='/project/{project.name}/EditorArea'>{project.name}</a>
                </li> -->
            {/each}
            </ul>
        {/if}

        
        <button class="button-b outfit-300" onclick={goToAllDocumentsPage}>Show all my projects</button>
        
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

    .MyProjects {
        height: 800px;
        width: 47%;
        float: left;
        margin: auto;
        padding-left: 2rem;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        align-items:baseline;
        background: rgb(250, 250, 250);
    }

    .title-container{
        display: flex;
        flex-direction: row;
        align-items: center;
        width: 90%;
        justify-content:space-between;
        margin-top: 1rem;
        white-space: nowrap;

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

    .icon {
        font-size: 20px;
        font-weight: bold;
        }

    #create-new-project-popup-overlay {
        position: fixed;
        width: 100%;
        height: 100%;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: rgba(0,0,0,0.5);
        z-index: 2;
    }

    #create-new-project-popup {
        background-color: white;
        position: fixed;
        top: 25%;
        left: 50%;
        transform: translate(-50%, -50%);
        box-shadow: 2px 3px 5px #999;
        width: 25%;
    }

    #create-new-project-popup > #header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin: 0 1rem;
    }

    #create-new-project-popup > #header > h1 {
        display: inline-block;
        font-size: 1.5em;
    }

    #create-new-project-popup > #header > #closeBtn {
        cursor: pointer;
    }

    #create-new-project-popup > #header > #closeBtn:hover {
        fill: #fa0000;
    }

    #create-new-project-popup > input {
        display: block;
        font-size: 1rem;
        margin: 0 1rem;
        width: calc(100% - 2rem);
        padding: 6px 8px;
        box-sizing: border-box;
        border: 1px solid #677283;
        border-radius: 4px;
        transition: border-color .15s ease-in-out;
        outline: none;
    }

    #create-new-project-popup > input:focus {
        border-color: #366cbf
    }

    #create-new-project-popup > #project-name-error-msg {
        margin: 0 1rem;
        color: red;
    }

    #create-new-project-popup > #submitBtn {
        background-color: #34a93c;
        color: #fff;
        margin: 0.5rem 1rem;
        padding: 10px 25px;
        border: 1px solid #34a93c;
        border-radius: 4px;
        box-shadow: rgba(0, 0, 0, .2) 0 2px 4px 0;
        font-size: 16px;
    }

    #create-new-project-popup > #submitBtn:enabled {
        cursor: pointer;
    }

    #create-new-project-popup > #submitBtn:hover:enabled {
        background-color: #168e48;
    }

    #create-new-project-popup > #submitBtn:disabled {
        background-color: #8FC493;
        border: none;
    }

    #projects-separator {
        border: none;
        border-top: 2px solid #393939;
        margin: 20px 0 0 0;
        width: 90%
    }

    ul {
        padding: 0;
    }

    .project-container{
        display: flex;
        flex-direction: row;
        background-color: rgb(255, 255, 255);
        margin-bottom: 1rem;
        align-items: center;
        padding: 1rem;
        border-radius: 4px;
        width: 80%;
        box-shadow: rgba(0, 0, 0, .2) 0 2px 4px 0;
        color: inherit;
        text-decoration: none; 

    }

    .project-container:hover{
        cursor: pointer;
        background-color: #e0e0e0;
    }

    .docIcon{
        width: 30px;
        height: 30px; 
        margin-right: 1rem;
    }

    .projects {
        flex-grow: 1;
        width: 100%;
        color: black;
        padding: 1rem;
    }

    .button-b {
        background-color: #34a93c;
        border: 1px solid #34a93c;
        border-radius: 4px;
        box-shadow: rgba(0, 0, 0, .2) 0 2px 4px 0;
        box-sizing: border-box;
        color: #fff;
        cursor: pointer;
        font-size: 16px;
        outline: none;
        outline: 0;
        padding: 10px 25px;
        text-align: center;
        user-select: none;
        touch-action: manipulation;
    }

    .button-b:hover{
        background-color: #168e48;

    }

</style>
