# goapp-tutorial

Companion code for the article
[Creating a scalable API with Go, Gin and MongoDB Atlas](https://ayaacodes.hashnode.dev/creating-a-scalable-api-with-go-gin-and-mongodb-atlas) (2023).

A small Gin API with sign-up, sign-in and a JWT-protected route group, backed by
MongoDB. It shows one way to lay out a Go service: `cmd/web` for the entrypoint
and routes, `handlers` for HTTP, and `modules/` for auth, config, database access
and password hashing.

```sh
go run ./cmd/web
```

Routes: `GET /`, `POST /sign-up`, `POST /sign-in`, and `GET /auth/dashboard`
behind the `Authorization` middleware. Set the MongoDB connection string in your
own environment before running.

Archived. Kept for readers of the article; not maintained.
