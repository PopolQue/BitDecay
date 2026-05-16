# Game Design Document – _Campers_

## Executive Summary
**Campers** ist ein rundenbasiertes Online-Multiplayer mit strategischen und wirtschaftlichen Mechaniken. Die Spieler:innen übernehmen die Leitung eines Campingplatzes, verwalten Ressourcen, bauen Schlafplätze aus, managen Gäste und konkurrieren um die meisten Siegespunkte am Ende einer Saison.

- **Genre:** Strategie, Simulation, Deck-Building 
- **Spieleranzahl:** 2–8
- **Dauer:** 70–150 Minuten
- **Alter:** 14+
- **Projektumfang:** Webgame mit Karten-, Token- und Ressourcenmanagement.
- **Ziel:** Entwicklung eines prototypischen Core Games mit Erweiterungspotenzial für Spin-offs (z. B. Familienversion, digitale Umsetzung).
    
## Core Concept
### Concept Statement
„Baue, erweitere und manage deinen Campingplatz, um Gäste zufrieden zu stellen, Einnahmen zu generieren und am Ende die meisten Siegespunkte zu sammeln.“    

### Player Experience and Game POV
- Ziel: Am Ende des Spiels die meisten Siegespunkt haben
- Spieler:innen erleben Spannung durch Ressourcenknappheit und Konkurrenz
- Humorvolle Identifikation mit Gästen & Situationen
- Mischung aus strategischer Planung und taktischen Reaktionen auf gezogene NPC-Karten

## Main Features

### Story
Das Spieljahr (12 Runden) wird in 4 Quartale à 3 Runden unterteilt 

| Quartal | Saison      | Ereignisse   | Gästeaufkommen           |
| ------- | ----------- | ------------ | ------------------------ |
| Q1      | Nebensaison | 1 Ereignis   | -1 NPC-Karte pro Spieler |
| Q2      | Hauptsaison | 2 Ereignisse | +1 NPC-Karte pro Spieler |
| Q3      | Hauptsaison | 2 Ereignisse | +1 NPC-Karte pro Spieler |
| Q4      | Nebensaison | 1 Ereignis   | -1 NPC-Karte pro Spieler |

Dramaturgischer Verlauf: Aufbau → Boom → Peak → Abkühlung

Dieses System erzeugt:
- Strategiewechsel im Spielverlauf
- Zyklische Ressourcenknappheit
- Story-Momente am Tisch

#### Ereigniskarten-System
Grundprinzip: Ereignisse werden zu Beginn eines Quartals aufgedeckt
- Sie gelten für 3 Runden
- Sie betreffen alle Spieler gleichzeitig
- Sie verändern Nachfrage, Ressourcen oder Asset-Werte

Modularer Aufbau:
- X Basiskarten
- Y Wetterkarten-Modul
- Z Wirtschafts-/Trendkarten-Modul
- W High-Interaction-Karten (Abstimmungen)

##### Entwicklungs-Ereignisse // (siehe 0226_Erweiterung für mehr aber erstmal ignorieren) 
- Dauerregen:
    - Es regnet in Strömen und für diese Saison ist kein Ende in sicht..
    Effekt:
        - NPCs pro Zug halbieren sich
- Tourismusboom:
    - Eure Region wurde zum Weltkulturerbe ausgerufen!
    Effekt:
        - NPCs pro Zug verdoppeln sich
- Sturmwarnung:
    - Das Aufeinandertreffen von Hoch- & Tiefdruckgebieten sorgt für heftige Unwetter in den Nächten.
    Effekt:
        - NPCs wollen nicht in Zelten untergebracht werden
- Festival:
    - In eurer Region findet ein Festival statt
    Effekt:
        - NPCs wollen bevorzugt in Zelten untergebracht werden
- Hitzewelle:
    - Die pralle Sonne macht eure Besuchenden extra durstig!
    Effekt:
        - erhöhter Wasserverbrauch
- Dürre:
    - Der Nestel-Konzern hat so viel Grundwasser abgepumpt, dass es zu einer Wasserknappheit in eurer Region gekommen ist.. 
    Effekt:
        - Wassertanks generieren kein Wasser/Runde
- Blackout:
    - Ein Sonnensturm hat das lokale Energienetzwerk flachgelegt..
    Effekt:
        - Stromgenerator generieren keinen Strom/Runde

### Gameplay
- Rundenbasiertes Ressourcen- und Gäste-Management.

#### Core Game Loops
1. Am Anfang eines Zuges wird der NPC-Pool aufgefüllt (potenzielle Gäste)
2. Führe bis zu X Aktionen aus:
    - **Bauen**: Assets errichten oder upgraden
    - **Kaufen**: Ressourcen erwerben
    - **Gästemanagement**: Buchungen annehmen, Gäste platzieren oder entfernen
3. Am Rundenende: Einkommen durch aktive Gäste erhalten

<img width="894" height="680" alt="Bildschirmfoto 2026-02-24 um 18 47 44" src="https://github.com/user-attachments/assets/d4c6deb0-9745-4553-b04e-936323af41ea" />

#### Game Feature

![Campers - Frame 2](https://github.com/user-attachments/assets/01dcb53a-baee-4376-9b88-0d3e48cdfda2)

##### NPCs:
mandatory:
- Namen
- Typ // wie bisher (Hippies, Familien & Snobs)
- angefragter Zeitraum	
- Anzahl Gäste
- Einkommen/Nacht
- Einkommen gesamt
- Bedürfnisse (pro Nacht):
-- Anzahl Schlafplätze
-- Strombedarf
-- Wasserbedarf
- Sonderbedürfnisse:
-- bestimmte Assets müssen auf dem Zeltplatz vorhanden sein

##### Assets:
Assets als Interface mit den Werten:
- Preis
- Flächenbedarf
3 Arten von Assets:
- Buchbare Assets:
	- Schlafplätze:
        - Zelt (2 Personen) (upgradebar auf Glamping-Zelt (4 Personen))
        - Caravan (4 Personen)
        - Bungalow (6 Personen) (upgradebar auf Luxus-Bungalow (8 Personen) )
- Ressourcen generierende Assets:
    - Bisher Stromgenerator & Wassertank
- Sonderbedürfnisse erfüllende Assets:
    - Sportplatz (Sonderbedürfnis mancher Familien)
    - Lagerfeuerstelle (Sonderbedürfnis mancher Hippies & Familien)
    - Sauna (Sonderbedürfnis mancher Snobs)
    - Open-Air Bühne (Sonderbedürfnis mancher Hippies & Snobs)

###### NPC-Schlafplatz-Logik:
**Hippies** können Anfragen stellen für: Zelte (nicht mehr als 2 Personen), Caravans (nicht mehr als 4 Personen) oder Bungalows (nicht mehr als 6 Personen).
**Familien** können Anfragen stellen für: Glamping-Zelte (nicht mehr als 4 Personen), Caravans (nicht mehr als 4 Personen) oder Bungalows (nicht mehr als 6 Personen).
**Snobs** können Anfragen stellen für: Glamping-Zelte (nicht mehr als 2 Personen), Bungalows (nicht mehr als 4 Personen) oder Luxus-Bungalows (nicht mehr als 6 Personen).

##### Ressourcen:
Aufgeteilt in kaufbare (Fläche) & generierbare (Geld, Strom & Wasser)

Spieler können:
- Mit ausreichend Geld: 	
    - eine Fläche kaufen und die gesamt Fläche ihres Campingplatzes vergrößern.
- Mit ausreichend Geld & Fläche: 
    - Einen der 3 Asset typen kaufen & platzieren.
- NPC Anfragen sehen & diese, wenn sie die Bedürfnisse erfüllen können annehmen & einem den Bedürfnissen entsprechenden Schlafplatz zuweisen. Oder sie ablehnen.

## Game Menu

### Landing Page

Ein Homescreen, der das Campers-Logo zeigt und das Campers-Thema für die Sommersaison als Hintergrund hat.

User können Auswählen zwischen:

- Tutorial spielen
- Solo-Spiel spielen
- Multiplayer spielen
    - Lobby erstellen
    - Lobby beitreten
- Bestenliste anschauen


### Tutorial

User spielen 2 Saisons, in denen sichergestellt ist, dass alle Spielmechanismen auftretten und erhalten erklärende Pop-Ups jedes mal, wenn ein Spielmechanismus zum ersten Mal erscheint. Ebenso werden alle Resourcen- & sonstige Entitätenarten (NPC-Typen, Schlafplatz-Type, Asset-Typen) erklärt.

### Solo-Spiel

User spielen ein Spiel alleine und können anschließend ihren Score auf dem Leaderboard posten.

### Mulitplayer

Bis zu 8 User können in einer Lobby gemeinsam gegeneiander Spielen. Der Gewinner des Spiels kann anschließend seinen Score auf dem Leaderboard posten.

#### Lobby erstellen

User erstellen eine Lobby und generieren einen Beitrittscode, den sie an die anderen User versenden können.
    
#### Lobby beitretten

User können einen Beitrittscode eingeben und so einer erstellten Lobby beitretten.