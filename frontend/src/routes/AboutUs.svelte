<script>
    let Message = "TexMex is the fastest, easiest way to collaborate on LaTeX—completely free. Multiple users, zero lag, no installs. Just pure, seamless, real-time editing with all the power of LaTeX and none of the usual friction. It’s how LaTeX was meant to be done."
</script>

<div class="AboutUs">
   <div class="Header cal-sans-regular">
        <h1>About TexMeX</h1>
    </div>
    <div class="Textbox outfit-300">
        <p>{Message}</p>
    </div>

    <img src="/LogoTest.png" alt="Our Beautiful Logo" />

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

    .AboutUs {
        height: 800px;
        width: 49%;
        float: right;
        margin: auto;
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        background: rgb(250, 250, 250);

    }

    .Header{
        align-self: center;
    }

    .Textbox{
        width: 90%;
        margin-left: auto;
        margin-right: auto;
        word-wrap: normal;
    }

    img {
        display: block;
        margin: 0 auto;
 
        height: auto;
    }

</style>
