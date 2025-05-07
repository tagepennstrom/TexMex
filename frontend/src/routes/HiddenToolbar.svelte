<script lang="ts">
    //OBS!
    //Denna komponent gör inget, används för att homescreen headern skall vara lika stor som editor headern oavsett skärm, finns säker
    //bättre lösning, men tillsvidare funkar detta.
</script>

<div class="Toolbar">
    <button class="hidden">hidden</button>
</div>

<style>

    .Toolbar{
        display: flex;
        gap: 10px;
        padding: 10px;
        background-color: #e0f7fa;
        border-bottom: 1px solid #ccc;
        align-items: center;
    }

    .hidden{
        padding: 10px 20px;
        background-color: rgb(246, 137, 4);
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        visibility: hidden;
    }

</style>
