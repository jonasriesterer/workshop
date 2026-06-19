# Programmierworkshop am 19.6.2026

## Namen
Jonas Riesterer, Nikolas Vix

## Link zum Git-Repository
https://github.com/jonasriesterer/workshop

## KI-Werkzeuge
Devin Desktop, Claude

### Agenten
Cascade, Claude Code

### Chat-URLs
https://claude.ai/

## Frameworks und Bibliotheken
Gin

### REST-Schnittstelle (Lesen und Neuanlegen)

### Validierung (nur Neuanlegen)

### OR-Mapping (für PostgreSQL)

### Optional: OIDC mit Keycloak

### Einfacher Integrationstest

## Prompts/Requests an KI-Agent/en

### Claude
- Ich brauche in diesem Ordner die Struktur für einen Server in Go der auf dem Framework Gin basiert. erstell mir die relevante struktur und dateien für eine laufende server anwendung ohne Logik
- Ich habe die struktur jetzt übernommen. jetzt gib mir eine liste von wichtigen commands zum lokalen server start/stop, dev modus, tests, usw... gerne auch allgemeine gin commands für validierung etc
- Der Server läuft jetzt mir air. Jetzt will ich anfangen, die Funktionen zu implementieren. Ich will als erstes 3 Entities im Ordner internal/model erstellen. die Hauptentität Autohaus als aggregate root, über die alle zugriffe laufen, dann als 1 - n beziehung Autohaus - Auto und als 1 - 1 beziehung Autohaus - Adresse. erstell mir die notwendigen Dateien so, dass ich danach die Eigenschaften der Entities eintragen kann
- Mein Autohaus soll folgende Eigenschaften haben:
id, name, username, email, AnzahlFahrzeuge, gründungsdatum, homepage, telefonnr, version, erzeugt und aktualisiert. füge diese eigenschaften dem model hinzu mit passenden Kommentaren und Validierungen
- jetzt mach das gleiche für meine entities auto und adresse. für Auto:
id, kennzeichen, marke, modell, baujahr, autohaus
für Adresse:
id, plz, ort, land, Autohaus
- ich habe sql dateien zum erstellen von den Datenbanktabellen, indizes, schema usw. wo lege ich diese dateien nach Go standard ab? außerdem habe ich noch sql dateien, die mir direkt testdaten zu den entities generieren. kommen die an den gleichen ort oder an einen anderen?
- ich will noch ein kleines frontend bauen mit dem ich per knopf requests schicken kann. wie baue ich das in meine struktur ein und was brauche ich alles dafür um es in go umzusetzen
- funktioniert, aber es sieht nicht schön aus. ich will die response nicht als json datensatz, sondern die einzelnen eigenschaften immer mit einem passenden icon angezeigt bekommen untereinander
- ich will jetzt noch weitere endpunkte einbauen. einmal standard put, einmal standard delete. dann will ich noch bedingte get requests haben mit 304 not modified als antwort, wenn die versionsnummer sich nicht verändert hat. bei dem update mit put soll kein body zurück gegeben werden
- ich brauche noch ein docker file, damit ich meinen server auch in docker laufen lassen kann. erstell mir das nötige file
