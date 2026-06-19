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

Gin (`github.com/gin-gonic/gin v1.10.0`)

### Validierung (nur Neuanlegen)

go-playground/validator mit `validate:"..."` Struct-Tags (z. B. `required`, `email`, `min`, `max`, `url`, `e164`)

### OR-Mapping (für PostgreSQL)

SQLC zur Code-Generierung aus SQL-Queries + pgx als PostgreSQL-Treiber

### OIDC mit Keycloak

JWT Bearer Token Validierung über Keycloak OIDC Discovery (JWKS).

- Realm: `go_workshop`, Client: `go-client`
- Rollen (Client-Rollen): `admin`, `user`
- `POST /rest`, `PUT /rest/:id` → Rolle `admin` oder `user` erforderlich (sonst `401`/`403`)
- `DELETE /rest/:id` → nur Rolle `admin` (sonst `401`/`403`)
- `GET /rest/:id` → öffentlich, kein Token nötig
- TLS-Verify für self-signed Keycloak-Zertifikat deaktivierbar via `KC_SKIP_TLS_VERIFY=true`

### Einfacher Integrationstest

testcontainers-go startet einen echten PostgreSQL-Container; Assertions mit testify

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

### Cascase (Devin Desktop)

Prompt 1:

Ausgangssituation

* Go als Plattform
* DB-Server aus den vorangegangenen Abgaben
* Keycloak ist optional


Wir müssen einen Backend server in Go implementieren, kannst du mir einen Plan machen wie das ganze aufgesetzt werden muss. Ich hab bisschen recherchiert und mich für folgende tools entschieden:
Sprache? Go
Web-Framework? Gin
Linter? golangci-lint
Formatter? goimports
Logger? log/slog
ORM? sqlc

Testing?
- testing:  test.go datei endungen --> mit go test laufen lassen, unterstützt unit tests, benchmarks, fuzzing
-stretchr/testify
-testcontainers-go für integrationstests 


Kannst du das validieren und darauf basierend dann einen plan erstellen wie die ganzen tools in das bereits bestehende Setup integriert werden können?

Wir haben ein datenbank setup (sql dateien) und entitäten von einem anderen Projekt, die wir verwenden sollen, da kannst du dann auch mal schauen wie wir das integriert bekommen. Es soll einen Test-Modus geben (Umgebungsvariable), wo bei "true" die Datenbank neulädt und die csv daten reinlädt (DB Seeding) für den Testmodus, bei "false" soll das ganze nicht passieren.

Die Entitäten existieren bereits
Implementiere eine Schnittstelle mit  ersteinmal einen GET Request und einem POST Request passend zur Fachdomäne. Der Zugriff soll über das AggregateRoot "Authohaus" geschehen. Nutze dabei die bestehende Struktur mit handler, model (hier sind die entitäten), service für Geschäftslogik und repository für den DB-Zugriff. Nutze für POST ein sinnvolles Validierungstool

Implementiere dann die Tests (Unit- und Integrationstests), ein paar wenige reichen.

Verwende bitte alle oben gennanten Tools sinnvoll.

Prompt 2:

Kannst du noch keycloak für die POST, PUT und DELETE Endpunkte security mit keycloak einbauen, in der readme unter keycloak sieht man, dass es nur die Rolle admin und user gibt, bei delete braucht man die rolle admin ansonsten passt admin oder user.


## Projektbeschreibung

Eine REST-API für ein Autohaus-Verwaltungssystem, gebaut mit Go und dem Gin-Framework.

### Endpunkte

| Methode | Pfad | Beschreibung |
| --- | --- | --- |
| `GET` | `/rest/:id` | Autohaus abrufen (mit `ETag` / `If-None-Match` → `304 Not Modified`) |
| `POST` | `/rest` | Neues Autohaus anlegen (→ `201 Created` mit `Location`-Header) |
| `PUT` | `/rest/:id` | Autohaus aktualisieren (`If-Match` für optimistisches Locking, kein Body → `204`) |
| `DELETE` | `/rest/:id` | Autohaus löschen (→ `204 No Content`) |
| `POST` | `/dev/db_populate` | Datenbank zurücksetzen und neu befüllen |
| `GET` | `/health` | Health-Check |

### Web-Oberfläche

Unter `http://localhost:3000` ist eine einfache Web-Oberfläche erreichbar, über die GET-Requests per Knopfdruck abgeschickt werden können.

### Technologien

- **PostgreSQL** — Datenbank mit SQLC-generiertem Datenbankzugriff, optimistisches Locking über Versionsnummer
- **Keycloak** — OIDC-Authentifizierung (vorbereitet)
- **Docker** — Der Server läuft als Container, gestartet über Docker Compose zusammen mit Postgres und Keycloak

### Docker

```bash
# Alle Services starten (Autohaus-Server, Postgres, Keycloak)
docker compose -f extras/compose/autohaus/compose.yml up --build -d

# Stoppen
docker compose -f extras/compose/autohaus/compose.yml down
```

### Nützliche Commands

```bash
# Server lokal starten
go run ./cmd/server
air                                  # mit Live-Reload

# Tests
go test ./...                        # alle Tests

# Build
go build ./...                       # kompilieren (Fehlerprüfung)

# Validierung & Formatierung
go vet ./...                         # statische Analyse
go fmt ./...                         # Code formatieren
golangci-lint run ./...              # Linter
```