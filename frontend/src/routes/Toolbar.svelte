<script lang="ts">
    import { insertBold } from '$lib/index';

    function makeBold() {
        insertBold();
    }
</script>

<button class="Bald" onclick={makeBold}>Bold</button>


<style>
    .Bald {
        padding: 10px 20px;
        background-color: rgb(246, 137, 4);
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }
</style>
