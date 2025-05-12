<script lang='ts'>
	import { toggleFilesModal } from '$lib';
    import { showFilesModal } from './stores';
    import {fly } from 'svelte/transition';
</script>

{#if $showFilesModal}
<div class="backdrop" transition:fly="{{x: -1000, duration: 300}}">
    <div class="modal">
        <button class="close" onclick={toggleFilesModal}>Close</button>
        <input type="file" id="file" />
        <label for="file">Upload File</label>
    </div>
</div>
{/if}

<style>
    .backdrop{
        width: 50%;
        height: 50vh;
        background-color: #e0f7fa;
    }

    .modal {
        border-radius: 10px;
    }

    Input {
        Display: none;
    }

    Label {
        Height:  100px;
        Width: 150px;
        Border-radius: 60px;
        Border: 1px dashed #999;
    }

    Label:hover {
        Color: #de0611;
        Border: 1px dashed #de0611;
    }
</style>


