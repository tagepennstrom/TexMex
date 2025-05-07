<!---- THIS COMPONENT IS NOT USED AS OF NOW, WE USE THE HEADER IN THE ROUTES TAB>
<script>
    let menuOpen = false;
</script>

<div class="Header">
    <img src="/LogoTest.png" alt="Our Beautiful Logo" />
    <button class="menu-button" on:click={() => menuOpen = !menuOpen}>
        ☰
    </button>
        
        {#if menuOpen}
        <nav class="dropdown">
            <a href ="/">Home</a>
            <a href ="/EditorArea">Editor</a>
        </nav>
    {/if}
</div>

<style>

    .menu-button{
        position: absolute;
        top: 1rem;
        left: 1rem;
        background: none;
        border: none;
        font-size: 1.5rem;
        cursor: pointer;
    }

    .dropdown{
        position: absolute;
        top: 3rem;
        left: 1.5rem;
        background: white;
        border: 1px solid #ccc;
        padding: 10px;
        font-size: 1 rem;
        cursor: pointer;
    }

    .Header{
        background-color :#e0f7fa;
    }

    img {
        display: block;
        margin: 0 auto;
        max-width: 10%;
        background-color :#e0f7fa;
        height: auto;
    }
</style>