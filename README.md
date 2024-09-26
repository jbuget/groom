# Installation

Démarrer Postgres via Docker Compose

```shell
docker compose up -d
```

Configurer l'environnement

```shell
export DATABASE_URL=postgres://postgres:password@localhost:5432/ushr \
    BASIC_AUTH_LOGIN=admin \
    BASIC_AUTH_PASSWORD=admin \
    USHR_API_KEY=your_api_key_here
```

Initialiser le projet

```shell
go mod init ushr
```

To add module requirements and sums:

```shell
go mod tidy
```

Lancer l'application

```shell
go run ./cmd/ushr
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
curl -X POST http://localhost:3000/api/rooms -d '{"slug":"nouvelle-salle", "meet_id":"xxx-yyyy-zzzz"}' -H "Content-Type: application/json" -H "X-API-KEY: your_api_key_here" 

# Modifiez une room existante
curl -X PUT http://localhost:3000/api/rooms/2 -d '{"slug":"salle-existante", "meet_id":"xxx-yyyy-zzzz"}' -H "Content-Type: application/json" -H "X-API-KEY: your_api_key_here" 

# Supprimez une room
curl -X DELETE http://localhost:3000/api/rooms/1
```