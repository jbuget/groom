# Groom, the Google rooms manager

## Installation

### Requirements

```shell
brew install golang-migrate
```

Générer des identifiants de type "compte de service" depuis la console Google Cloud.
Récupérer le fichier de credentials, `account_service.json`, à mettre à la racine.

### Steps

Démarrer Postgres via Docker Compose

```shell
docker compose up -d
```

Configurer l'environnement

```shell
export DATABASE_URL="postgres://postgres:password@localhost:5432/groom?sslmode=disable"
export GROOM_API_KEY="your_api_key_here"
export GOOGLE_API_KEY="<your_google_api_key>"
export GOOGLE_CLIENT_ID="<your_google_client_id>"
export GOOGLE_CLIENT_SECRET="<your_google_client_secret>"
export GOOGLE_REDIRECT_URL="https://example.test/auth/callback"
export GOOGLE_WORKSPACE_DOMAIN="example.test"
export GOOGLE_SERVICE_ACCOUNT_IMPERSONATED_USER="service-account@example.test"
export GOOGLE_SERVICE_ACCOUNT_CREDENTIALS_FILE="./service_account.json"
```

Initialiser le projet

```shell
go mod init groom
```

To add module requirements and sums:

```shell
go mod tidy
```

Lancer l'application

```shell
go run ./cmd/groom
````

# Usage

## URL HTML

```shell
# Afficher la liste des rooms
http://localhost:3000

# Accéder à une room 
http://localhost:3000/ma-room
```

## API

```shell
# Lister les rooms
curl http://localhost:3000/api/rooms -H "X-API-KEY: your_api_key_here" 

# Ajouter une room
curl -X POST http://localhost:3000/api/rooms -d '{"slug":"nouvelle-salle"}' -H "Content-Type: application/json" -H "X-API-KEY: your_api_key_here" 

# Modifiez une room existante
curl -X PUT http://localhost:3000/api/rooms/2 -d '{"slug":"salle-existante", "space_id":"xxx-yyyy-zzz"}' -H "Content-Type: application/json" -H "X-API-KEY: your_api_key_here" 

# Supprimez une room
curl -X DELETE http://localhost:3000/api/rooms/1 -H "X-API-KEY: your_api_key_here" 
```


## How to

### Manually excute migrations

```shell
migrate -path ./migrations -database ${DATABASE_URL} up
migrate -path ./migrations -database ${DATABASE_URL} down [1]
```
