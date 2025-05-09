<script lang='ts'>
	import AboutUs from './AboutUs.svelte';
	import Footer from './Footer.svelte';
    import Header from './Header.svelte';
    import HiddenToolbar from './HiddenToolbar.svelte';
	import MyProjects from './MyProjects.svelte';

</script>



<div class="page-container">
    <Header/>
    <div class="hiddenToolbar">
        <HiddenToolbar/>
    </div>
    <div class ="contents">
        <MyProjects/>
        <AboutUs/>
    </div>
    <Footer/>
</div>

<style>
    .page-container {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
        background-color: rgb(250, 250, 250);
    }

    .contents{
        
    }


    .Footer{
        align-self: center;
    }


</style>