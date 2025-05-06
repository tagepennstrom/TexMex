<script lang='ts'>
    import Header from './Header.svelte';
    import HiddenToolbar from './HiddenToolbar.svelte';

    function handleClick() {
        window.location.href = "/EditorArea";
    }
</script>



<div class="page-container">
    <Header/>
    <div class="hiddenToolbar">
        <HiddenToolbar/>
    </div>
    <div class="content">
        <button class="ToEditor" onclick={handleClick} title ="ToEditor"><strong>To the editor</strong></button>
    </div>
</div>

<style>
    .page-container {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }

    .ToEditor{
        position: relative;
        top: 1rem;
        left: 1rem;
        background: none;
        border: none;
        font-size: 1.5rem;
        cursor: pointer;
    }

    .content {
        flex: 1;
    }
</style>