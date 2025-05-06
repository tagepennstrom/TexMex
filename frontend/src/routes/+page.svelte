<script lang='ts'>
    import Header from './Header.svelte';

    function handleClick() {
        window.location.href = "/EditorArea";
    }
</script>



<div class="page-container">
    <Header/>
    <div class="content">
        <button class="Bold" onclick={handleClick} title ="ToEditor"><strong>To the editor</strong></button>
    </div>
</div>

<style>
    .page-container {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }

    .content {
        flex: 1;
    }
</style>