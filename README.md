# Golang-notes-app

notes-app backend in golang. The app can run in three APP_MODE: [Dev, Test, Prod]. 
As default, it will run in `dev` mode. When you run test, it will switch to run in `test` mode.
`Prod` is a bit strict running environment, you need to define the required environment variables before running the main app.

See `./config/config.go` to see how to define those environments. 

It is easy to extend the APP_MODE and add your preferred config, e.g. [Staging].      

# Development

## Requirements

In order to test develop locally, you need:

- Install `go`.
- Install `docker` for running postgres. 
- Add an `DEVROOT` in your `.bashrc` to define where to put the `dev` postgres volume data. See `Makefile` and search `DEVROOT`. 

After you have installed the above requirements. You are good to go.

## Useful commands

- `make start` will start the `main.go` app in dev mode.
- `make db-start` will start only the postgres docker.
- `make db-stop` will stop the postgres docker.
- `make db-ssh` will login the postgres docker container.

# Test

Run `make test` will run all the unit tests by `ginkgo`.

In the `test` mode, each test will spin up a new TestApp that forks a new `database` with a random database name and a new `port`.
After the whole test suite finishes, the database session pools will cleanup all the forked new database.
see `./db/sessionpool.go` to understand how it works. Each test will start from a fresh new database.

The benefit is that we can run all the tests independently in parallel without interfering with each other.    

# Production

## Deployment