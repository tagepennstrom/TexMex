<script>
    let menuOpen = false;
</script>

<div class="Header">
    <button class="menu-button" on:click={() => menuOpen = !menuOpen}>
        ☰
    </button>

    <svg class="tm-logo" xmlns="http://www.w3.org/2000/svg" version="1.0" viewBox="0 0 1000.000000 1000.000000" preserveAspectRatio="xMidYMid meet">

        <g transform="translate(0.000000,1000.000000) scale(0.100000,-0.100000)" fill="#000000" stroke="none">
        <path d="M6071 8469 c-47 -44 -51 -117 -8 -157 12 -12 54 -31 92 -43 39 -12 165 -71 280 -132 157 -82 224 -123 265 -162 81 -77 158 -135 215 -160 27 -12 81 -50 119 -84 38 -33 73 -61 78 -61 5 0 24 -11 42 -25 47 -36 167 -83 288 -114 103 -26 109 -29 185 -94 43 -36 129 -110 191 -162 62 -53 130 -120 151 -148 111 -146 128 -170 156 -217 17 -28 44 -73 60 -100 17 -26 39 -63 50 -82 11 -18 31 -51 45 -73 14 -22 40 -65 58 -95 18 -30 50 -81 70 -112 34 -52 54 -78 162 -213 77 -96 133 -185 144 -228 3 -14 11 -39 18 -54 18 -46 57 -153 79 -223 11 -36 27 -83 34 -105 13 -40 18 -61 50 -225 9 -47 20 -98 25 -115 4 -16 11 -48 15 -70 3 -22 12 -58 19 -80 14 -41 37 -123 61 -210 7 -27 23 -84 34 -125 12 -41 26 -91 30 -110 20 -81 34 -137 41 -170 17 -72 41 -177 56 -240 8 -36 19 -83 23 -105 4 -22 13 -58 19 -80 39 -145 52 -301 52 -638 0 -568 -34 -657 -280 -736 -187 -60 -294 -79 -485 -82 l-90 -2 -7 74 c-3 41 -11 103 -17 139 -7 36 -16 91 -21 122 -4 31 -13 72 -18 90 -6 18 -15 51 -20 73 -19 74 -60 176 -85 208 -14 18 -44 69 -68 112 -156 290 -148 270 -152 372 -3 98 -26 152 -84 204 -16 14 -34 43 -40 65 -15 53 -34 130 -44 182 -5 24 -20 65 -34 90 -30 57 -65 136 -95 217 -12 33 -35 89 -51 125 -15 36 -37 92 -48 125 -22 66 -57 160 -91 245 -12 30 -28 73 -36 95 -18 54 -55 153 -93 248 -17 43 -31 85 -31 92 0 8 -27 45 -61 82 -64 71 -110 150 -128 223 -16 61 -108 257 -141 300 -15 19 -35 50 -44 67 -9 18 -37 56 -62 85 -26 28 -53 68 -61 87 -8 20 -32 59 -53 86 -75 97 -170 275 -170 316 0 22 -111 198 -139 221 -56 46 -297 261 -321 288 -14 15 -41 38 -61 52 -41 27 -120 92 -159 128 -14 14 -48 44 -75 66 -28 23 -64 55 -81 70 -67 64 -236 172 -304 195 -19 6 -48 17 -65 24 -60 27 -170 68 -200 76 -47 12 -112 38 -95 39 35 1 -11 21 -93 40 -48 11 -105 29 -127 40 -62 31 -157 50 -260 50 -52 0 -124 7 -160 16 -48 11 -119 16 -265 16 -132 0 -221 5 -262 14 -70 15 -234 5 -353 -22 -127 -28 -156 -34 -220 -43 -36 -5 -105 -8 -155 -6 -219 6 -576 -53 -745 -122 -52 -22 -131 -54 -175 -71 -97 -37 -196 -86 -250 -122 -22 -15 -51 -35 -65 -44 -38 -25 -106 -81 -238 -196 -67 -58 -142 -121 -167 -140 -104 -79 -408 -376 -475 -463 -22 -29 -60 -77 -85 -107 -24 -30 -56 -71 -70 -90 -14 -19 -46 -62 -70 -95 -24 -33 -63 -91 -86 -130 -59 -97 -219 -420 -219 -441 0 -4 -9 -28 -21 -53 -12 -25 -32 -68 -44 -96 -40 -86 -91 -252 -120 -390 -9 -44 -22 -105 -30 -135 -7 -30 -18 -82 -25 -115 -7 -33 -18 -80 -25 -105 -23 -80 -28 -169 -28 -490 1 -207 6 -379 17 -500 19 -212 26 -258 68 -441 5 -26 43 -111 84 -190 164 -318 175 -346 183 -433 9 -98 22 -125 111 -226 86 -97 106 -105 310 -125 85 -8 259 -26 385 -40 431 -45 720 -76 895 -95 96 -11 285 -31 420 -45 135 -14 324 -35 420 -45 96 -10 265 -28 375 -39 110 -11 236 -25 280 -31 44 -6 170 -20 280 -31 110 -11 279 -29 375 -39 96 -11 285 -31 420 -45 427 -45 690 -73 880 -94 102 -11 295 -32 430 -46 135 -14 324 -35 420 -45 96 -11 281 -30 410 -44 129 -14 303 -32 385 -41 124 -14 165 -15 235 -6 47 7 163 16 257 22 167 9 174 9 215 -12 33 -17 63 -22 142 -24 55 -1 116 2 135 7 20 5 59 13 87 18 28 4 60 12 70 18 11 5 36 12 57 15 41 6 51 11 267 115 135 66 170 81 260 112 107 37 220 105 290 174 39 39 79 84 87 101 9 16 26 44 39 61 22 28 69 147 80 204 3 14 13 57 22 95 15 60 17 131 17 505 0 486 0 495 -27 635 -29 154 -31 165 -74 335 -7 30 -20 87 -29 125 -8 39 -19 86 -24 105 -23 92 -34 138 -42 175 -4 22 -17 72 -29 110 -12 39 -28 95 -35 125 -25 100 -37 145 -50 180 -12 32 -24 84 -56 235 -6 28 -15 66 -20 86 -5 20 -9 48 -9 62 0 23 -12 71 -46 177 -8 25 -21 65 -28 90 -8 25 -26 77 -40 115 -15 39 -34 93 -42 120 -9 28 -19 55 -23 60 -4 6 -13 30 -20 55 -7 25 -26 63 -41 85 -15 22 -36 54 -46 71 -30 50 -85 123 -144 193 -31 36 -70 86 -86 111 -51 74 -104 159 -104 165 0 3 -15 27 -32 54 -18 27 -46 72 -63 101 -16 28 -43 73 -60 100 -16 26 -38 61 -47 77 -121 197 -270 364 -472 528 -54 44 -109 90 -123 103 -13 12 -47 40 -76 61 -111 85 -166 128 -213 169 -114 98 -159 132 -193 143 -20 7 -46 18 -57 26 -12 8 -56 25 -99 38 -42 13 -94 35 -115 49 -22 14 -42 26 -45 26 -3 0 -28 15 -55 34 -56 39 -73 42 -265 56 -77 5 -156 14 -175 20 -19 6 -91 13 -160 16 -69 4 -164 15 -212 25 -115 25 -158 24 -187 -2z m-1 -2281 c-5 -270 -78 -2659 -81 -2663 -7 -6 -643 46 -650 54 -3 3 0 294 7 646 21 1055 20 1006 8 992 -6 -7 -37 -69 -68 -140 -31 -70 -75 -166 -98 -215 -22 -48 -121 -265 -220 -482 -99 -217 -182 -398 -186 -401 -6 -5 -363 12 -476 23 -26 3 -32 12 -80 128 -29 69 -90 213 -135 320 -46 107 -116 274 -156 370 -40 96 -89 209 -108 250 l-35 75 -1 -723 -1 -724 -47 6 c-27 3 -138 13 -248 22 -110 8 -200 16 -201 17 -3 4 -22 2242 -20 2245 2 1 153 15 335 29 217 18 335 23 341 17 6 -6 42 -86 81 -179 179 -427 406 -971 452 -1082 28 -68 53 -123 56 -123 3 0 58 123 123 273 275 644 521 1212 525 1217 6 6 732 67 821 69 58 1 62 0 62 -21z m-2992 -464 l2 -242 -72 -6 c-40 -3 -145 -8 -233 -11 -88 -3 -189 -8 -225 -11 l-65 -6 16 -816 c9 -449 15 -819 13 -821 -4 -4 -437 28 -440 32 -1 1 -9 310 -18 687 -9 377 -19 734 -22 794 l-6 109 -147 -7 c-80 -3 -194 -9 -252 -12 l-106 -6 -8 216 c-4 120 -5 219 -3 221 7 7 1451 123 1508 122 l55 -2 3 -241z"/>
        </g>
    </svg>
        
    <div class="hidden">☰</div>


        {#if menuOpen}
        <nav class="dropdown outfit-400">
            <a class="drop-item" href ="/">Home</a>
            <a class="drop-item" href = "/SavedProjects">Saved Projects</a>
        </nav>
    {/if}
</div>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Outfit:wght@100..900&display=swap');

    .outfit-400 {
        font-family: "Outfit", sans-serif;
        font-optical-sizing: auto;
        font-weight: 400;
        font-style: normal;
    }

    .hidden{
        opacity: 0;
    }

    .Header{
        display: flex;
        padding: 1rem 0;
        flex-direction: row;
        background-color: #ceff6c;
        justify-content: space-between;
        box-shadow: 1px 5px 7px -5px rgba(0,0,0,0.15);
        z-index: 2;
    }

    .tm-logo{
        width: 70px;
        height: 70px;
        padding-right: 15px;
    }

    .menu-button{
        background: none;
        margin: 0 1.5rem;

        border: none;
        font-size: 2rem;
        cursor: pointer;
    }

    .dropdown, .drop-item {
        box-sizing: border-box;
    }

    .dropdown{
        display: flex;
        flex-direction: column;
        position: absolute;
        top: 5rem;
        left: 1.5rem;
        background: white;
        cursor: pointer;
        border: 1px solid black;
    }

    nav{
        font-size: 1rem;


    }

    .drop-item{
        display: flex;
        flex-direction: row;
        background-color: rgb(255, 255, 255);
        align-items: center;
        padding: 0.5rem 2rem;
        width: 100%;
        color: inherit;
        text-decoration: none; 
    }
    .drop-item:hover{

        background-color: rgb(231, 231, 231);
        
    }

    
</style>
