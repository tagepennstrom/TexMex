

## Overview
Ett dokument är en struktur som existerar hos varje lokal klient.
Den består av:
1. CursorPosition
2. Textcontent

Alla bokstäver sparas i Textcontent, och hela poängen med CRDTs är att vi vill att alla klienters textcontents ska se likadana ut.
Textcontent är en länkad lista, där varje nod beskriver ett tecken i texten.
Varje nod har metadata om vilken bokstav som är där, vem som skrev den, och vilken koordinat den har (detta är viktigt för synkronisering).
En textcontent måste alltid ha en tom head-nod.

## Terminologi
#### Coordinate (coord)
En koordinat är en int array. Den säger vart en bokstav ska ligga i förhållande till allt annat. 
Den kan läsas ungefär som en float. [2, 1, 1] -> 2.11

#### CursorPosition
Den noden i den länkade listan, där din keyboardcursor ligger i texten för tillfället.

#### Raise
Att man "ökar" längden på en koordinat, för att få plats med den mellan två andra.


## Alla (known) insertion-cases
Case 1: Insertion mellan två coords, plats finns mellan dom
Action: Endast en inkrementering behövs
Ex: Mellan [2] och [4] -> Ny koordinat: [3]

Case 2: Insertion mellan två coords, ingen plats finns mellan dom
Action: En raise krävs
Ex: Mellan [2] och [3] -> Ny koordinat: [2,1]

Case 3: Insertion mellan två coords, ingen plats finns mellan dom. En raise går inte, då koordinaten efter har redan tagit den raisen. 
Action: Raise med [0,1], eller så många nollor som krävs för att klämma sig emellan
Ex: Mellan [2] och [2,1] -> Ny koordinat: [2,0,1]

Case 4: Insertion på slutet av en textfil
Ex: Efter [5] -> Ny koordinat: [6]

Case 5: Insertion i början av en textfil
Action: Tekniskt sätt Case 2

Case X:
Ex: Mellan [2,1,1] och [2,2] -> Ny koordinat: [2,1,2]

Case Y:
Case Z:

