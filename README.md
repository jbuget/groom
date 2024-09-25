# Installation

Démarrer Postgres via Docker Compose

```shell
docker compose up -d
```

Configurer l'environnement

```shell
export DATABASE_URL=postgres://postgres:password@localhost:5432/ushr \
    BASIC_AUTH_LOGIN=admin \
    BASIC_AUTH_PASSWORD=admin
```

Initialiser le projet

```shell
go mod init ushr
```

To add module requirements and sums:

```shell
go mod tidy
```

Installer la dépendnce stdlib

```shell
go get github.com/jackc/pgx/v4/stdlib
```

Lancer l'application

```shell
go run main.go
````

# Usage

## URL HTML

```shell
# Afficher la liste des rooms
http://localhost:8080

# Accéder à une room 
http://localhost:8080/ma-room
```

## API

```shell
# Lister les rooms
curl -X GET http://localhost:8080/api/rooms

# Ajouter une room
curl -X POST http://localhost:8080/api/rooms -d '{"slug":"nouvelle-salle", "meet_id":"xxx-yyyy-zzzz"}' -H "Content-Type: application/json"

# Modifiez une room existante
curl -X PUT http://localhost:8080/api/rooms/2 -d '{"slug":"salle-existante", "meet_id":"xxx-yyyy-zzzz"}' -H "Content-Type: application/json"

# Supprimez une room
curl -X DELETE http://localhost:8080/api/rooms/1
```