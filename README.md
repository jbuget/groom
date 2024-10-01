# Groom, the Google rooms manager

## Installation

### Requirements

```shell
brew install golang-migrate
```

### Steps

Démarrer Postgres via Docker Compose

```shell
docker compose up -d
```

Configurer l'environnement

```shell
export DATABASE_URL=postgres://postgres:password@localhost:5432/groom?sslmode=disable \
    BASIC_AUTH_LOGIN=admin \
    BASIC_AUTH_PASSWORD=admin \
    GROOM_API_KEY=your_api_key_here \
    GOOGLE_API_KEY=<your_google_api_key>
    GOOGLE_CLIENT_ID=<your_google_client_id> \
    GOOGLE_CLIENT_SECRET=<your_google_client_secret> \
    GOOGLE_REDIRECT_URL=http://localhost:3000/auth/callback
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
curl -X POST http://localhost:3000/api/rooms -d '{"slug":"nouvelle-salle", "space_id":"xxx-yyyy-zzz"}' -H "Content-Type: application/json" -H "X-API-KEY: your_api_key_here" 

# Modifiez une room existante
curl -X PUT http://localhost:3000/api/rooms/2 -d '{"slug":"salle-existante", "space_id":"xxx-yyyy-zzz"}' -H "Content-Type: application/json" -H "X-API-KEY: your_api_key_here" 

# Supprimez une room
curl -X DELETE http://localhost:3000/api/rooms/1
```


## How to

### Manually excute migrations

```shell
migrate -path ./migrations -database ${DATABASE_URL} up
migrate -path ./migrations -database ${DATABASE_URL} down [1]
```
