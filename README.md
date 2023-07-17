# Hotel reservation backend

# Project environment variables(.env contents)
```
    HTTP_LISTEN_ADDRESS=:8080
    JWT_SECRET=something_super_secret
    MONGO_URI=mongodb://localhost:27017
    MONGO_DBNAME=hotel-reservation
    MONGO_TEST_DBNAME=hotel-reservation
```

## Project outline

- users -> book room from the hotel
- admins -> check reservations/bookings
- Authentication/Authentication -> JWT
- hotels -> CRUD API -> JSON
- rooms -> CRUD API -> JSON
- scripts -> database management -> seeding database, migrations...

## Resources

### MongoDB driver

Documentation

```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client

```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber

Documentation

```
https://gofiber.io
```

Installing gofiber

```
go get github.com/gofiber/fiber/v2
```

## Docker

### Installing mongodb as a Docker container

```
docker run --name mongodb -d mongo:latest -p 27017:27017
```
