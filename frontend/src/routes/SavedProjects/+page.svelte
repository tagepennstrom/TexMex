<script lang='ts'>
    import Header from '../Header.svelte';
    import HiddenToolbar from '../HiddenToolbar.svelte';

</script>



<div class="page-container">
    <Header/>
    <div class="hiddenToolbar">
        <HiddenToolbar/>
    </div>
    <div class ="contents">
        <p>Nothing to show right now</p>
    </div>
</div>

<style>
    .page-container {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }


</style>