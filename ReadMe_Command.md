# Go & Gin — Befehls-Referenz

## Server starten

```bash
# Direkt starten
go run ./cmd/server

# Gebuildetes Binary starten
go build -o bin/server ./cmd/server && ./bin/server

# Mit Umgebungsvariablen
PORT=9090 APP_ENV=production go run ./cmd/server

# Mit air (Live-Reload, empfohlen für dev)
go install github.com/air-verse/air@latest
air init
air
```

## Build

```bash
go build ./...                       # alles kompilieren (prüft auf Fehler)
go build -o bin/server ./cmd/server  # Binary erstellen
go clean                             # Build-Artefakte löschen
```

## Docker

```bash
# Image bauen
docker build -t autohaus .

# Mit Compose starten (inkl. Postgres & Keycloak, baut Image neu)
docker compose -f extras/compose/autohaus/compose.yml up --build

# Im Hintergrund starten
docker compose -f extras/compose/autohaus/compose.yml up --build -d

# Stoppen
docker compose -f extras/compose/autohaus/compose.yml down

# Logs anzeigen
docker compose -f extras/compose/autohaus/compose.yml logs -f autohaus
```

## Abhängigkeiten

```bash
go mod tidy                          # fehlende Pakete holen, ungenutzte entfernen
go get github.com/foo/bar            # neues Paket hinzufügen
go get github.com/foo/bar@v1.2.3     # bestimmte Version
go mod download                      # alle Abhängigkeiten cachen (z.B. in CI)
```

## Tests

```bash
go test ./...                        # alle Tests ausführen
go test ./internal/service/...       # nur ein Paket
go test -v ./...                     # verbose (zeigt jeden Test)
go test -run TestUserCreate ./...    # einzelnen Test nach Name
go test -cover ./...                 # mit Coverage-Ausgabe

# Coverage im Browser anzeigen
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## Code-Qualität

```bash
go vet ./...                         # statische Analyse (eingebaut)
go fmt ./...                         # Code formatieren
golangci-lint run ./...              # Linter
gofmt -l .                           # zeigt Dateien, die formatiert werden müssten

# Externe Tools (einmalig installieren)
goimports ./...      # wie gofmt + import-Sortierung
                     # go install golang.org/x/tools/cmd/goimports@latest

staticcheck ./...    # erweiterter Linter
                     # go install honnef.co/go/tools/cmd/staticcheck@latest
```
