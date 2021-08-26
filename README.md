# notes-app

This is a complete back-end app to demonstrate how to write REST API and GraphQL Server API. This app is production ready. There is another branch `nats-pubsub` that supports pubsub via nats to easily decouple and simplify the services.

notes-app backend in golang. The app can run in three APP_MODE: [Dev, Test, Prod].
As default, it will run in `dev` mode. When you run test, it will switch to run in `test` mode.
`Prod` is a bit strict running environment, you need to define the required environment variables before running the main app.

See `./config/config.go` to see how to define those environments.

It is easy to extend the APP_MODE and add your preferred config, e.g. [Staging].

You can access the [website](https://notes.biosave.org/)

## Requirements

In order to test develop locally, you need:

- Install `go`.
- Install `docker` for running postgres.

After you have installed the above requirements. You are good to go.

In order to simplify the docker build in different CI env and make sure we do not accidentally run another version, we use
`vendor` to manage all the dependencies.

After you have change the `go mod` dependencies, you could run `make sync-vendor` to keep in sync.

If you want to regenerate the proto code, you need to install `protoc` compiler, i.e. Protobuf Compiler.
Go to https://github.com/protocolbuffers/protobuf/releases.

Once you've unzipped the compiler archive you will have two directories: bin and include.
The contents of bin should be placed somewhere on your PATH. If you intend to use the included well known types then copy the contents of the include directory somewhere as well, for example into /usr/local/include/. Make sure that the contents of /usr/local/include and sub-directories are readable by your normal user.

## System

The development could be done on `Linux/MAC` and run all the `Makefile` commands with no extra effort.

You could also develop on Windows 10 Home/Pro.

- Develop on Windows 10's `wsl2` with `Docker Desktop` on Windows enabled with `wsl2`. Check `Docker Desktop`'s official
guide how to set this up. Then you can run everything and test on your wsl2's Linux (e.g. ubuntu 20.04). Make sure docker
is working on your wsl2. You could develop on your wsl2's linux file system, but DO NOT accessing Windows file system,
it is really slow with wsl2. Use vscode Remote wsl2 support, you could use your Windows `vscode` to connect to your remote
src code on wsl2.

- You do not like the wsl2 and vscode's remote wsl2, then you could develop on Windows native system with `GitBash`.
You won't easily get the `Makefile` work, but you do not need it if your develop using Jetbrain `Golang` IDE on Windows.
You could still run `ginkto` and `gqlgen` on Windows shell or Gitbash. Yes, you still need to start postgres docker on WSL2.
Then run `go build main.go` or in `Golang` IDE. Everything just works unless your can run `Makefile`. You still need to
run `go install xxxx` for your `tools.go` respectively.

## Useful commands

- `make start` will start the `main.go` app in dev mode.
- `make db-start-test` will start only the postgres docker.
- `make db-stop-test` will stop the postgres docker.
- `make db-ssh-test` will login the postgres docker container.

See the Makefile for more commands.

## Tools

- check `./tools.go` to see the entry to all tools we need in development process.

You could run `make tools` to install all the tools from the vendor folder.

Why we need a `tools.go` to management the tool dependencies?

Otherwise, `go mod tidy` will remove all of `tools` dependencies from the `go mod` and `vendor`.

## Authentication

### Dual-tokens policy will be used here to mitigate the JWT's drawback of `logout`.

- `access_token` will be used to access resources and has a short expired time duration. Once it is expired, the client
side need to periodically fetch an access token using the `refresh_token` before the `access_token`'s expiration time.
`access_token` will be stateless, we do not need to query the database to check the existence and validity.
- `refresh_token` will be used to access the `auth` end point to generate a new `access_token` and has a long expiration.
The `refresh_token` will not be state-less. It needs to query the database if the `refresh_token` is still valid,
when logout we need to set the `refresh_token` invalid or blacklisted.

### JWT encryption algorithm.

Two choices:

- `HS-256`. A simple one with only a secret key. Both token generation and verification use the same secret key, so-called
`symmetric`. If in your server side, you have control for both generation and verification on your own side, then it is
your choice and should be secure and good enough.
- `ES-256`. Use a key pair, private_key and public_key, so-called `non-symmetric`. The token generation will need your private key, and the verification
will use the public key to verify the signing and check the validity. This will be useful if you are holding the core service
and authenticates the token, but the verification could run on another party's code base, then it is better than use the public
to check the signing validity and token's validity. If you use the HS-256, you will have to let the other parties know
your 'secret' key, which could be a potential security leak for your core-service.

For most cases, it would be enough and easier go with `HS-256` to generate your JWT token.

## Database

Datastore is the important part for managing a microservice. `PostgresSQL` is a production ready and scalable choice.

Some key issues regarding testing and setup database store for production, development, test.

### ORM vs Raw SQL

We use `sqlx` mainly as an extension of `go`'s official `sql`. Their APIs are compatible and look very similar. I have
previous experience of using `ORM` (both in typescript's `TypeORM` and Python's `SQLAlchemy`) model to access the database where a lot of magic inside the `ORM` model and sometimes
it is really hard to know what happened when debugging and not easy to optimize when you hit on performance issues. On
the other hand, your complicated knowledge of learning a specific `ORM` library won't be transferred and reused for another
language or framework. Using the `vanilla` sql to query and `sqlx` to scan your expected returned queries are more
maintainable and scalable. You can easily switch to `nodejs` using typescript and other languages using raw sql style.

The choice here is to use the raw style to access the database. Even you switch to another database, it won't take too much
effort, absolutely the sql dialect could be an issue you have to fix.

### Database migration

It is a bit confusing that is called `migration` as terminology. Since when I first started, I thought it is only used for
`migration`, actually it is not. It involves every version change in sql schema from one version to another, e.g. your init
your database with init setup of some tables or fill in init data even the first time. It should be managed by the `migration`.
Then if you spin up a new database, it will go through all the migration versioned files and will keep track of in the database
itself some records to indicate if it has applied those migrations.

You can go to folder `./db/migrations` to check all the migration files. Any change to your database's schema should be
managed by the migration files.

In the file './db/migration.go', you could check how we apply those migrations for the database when starting up the app,
for example during the testing, we create a new database for every test, then it will apply all the migrations for the new
database. You could see the logging how it works.

### Test the data store and handlers

For testing, in the Makefile's `make db-start-test` we will start a docker without any volume for the `data`, with default
`postgres` username and password. When `make db-stop-test`, it will shutdown the docker and ofc everything will be trashed.

Instead of mocking the database for every test, so-called unit test, here we will spin up a new `database` with a random
name for every `TestApp`. The trick is that with one database connection to `postgres`, we could creat a new database.
We will use the new one for every single test. check the `./db/sessionpool.go` for more details how to fork a new database
base for `TestApp`. The `App` will use the `database` specified in the config, but the `TestApp` will use the database
specified in config for `Test` mode as a parent to spin up a new database for every test. That is the only difference 
between `TestApp` and `App`.


### Development locally with the database docker

Run `make db-start-dev`, it will start up a database docker for development only, it will use a volume `postgres-dev-notes-app`.
It is does not exist, it will create one. Docker volume is good and easy to use during dev process. Then next time you 
run the `dev` postgres, you will continue your last
database and data.

You can also run `make serve-prod` to spin a docker-compose for both postgres docker and app docker in production mode.


## HTTP Handler

### Golang web framework

Golang's built-in http framework is good enough and production-ready, normally it is preferred to build up your own with
some tiny middlewares or utilities. Here our choice is `chi` which simplifies a bit the golang's http handling without
losing the generality.

No matter if we use graphql, it is still valuable to write and handle the http requests in a good way. We use `chi` for
simplify our http handlers and http routers. It is more flexible and uses also standard `http.Handler` instead of using
specialized handlers, e.g. `gin` framework.

All the HTTP handlers are specified in `./handlers/handlers.go`. With utilities and middleware provided in `chi`, you gain
similar convenience as in `gin` but with more extra flexibility. It is also easy to interface with all other middlewares
not confined to `chi`'s middlewares, since `chi` is compatible with standard `http.Handler`.

### HTTP error handling

See the `./handlers/util/errors.go` and `./handler/util/render.go`, you can easily call `SendError` to send whatever you
want, it is also better you call those predefined errors in `./handlers/util/errors.go`.

### Payload and Response

For every http request, we should clearly define the `Payload` in a struct annotated with json and clearly define the `Resp`.
see examples in `todos.go` and `notes.go` inside `./handlers` folder. The `Payload` and `Response` may not exactly the
same as `models` defined in `./model`. The `model` is more general and communicates with `store` and other parts of the
application.

## GraphQL API

Compared to the `REST` API (i.e. `todos.go` and `notes.go` in `handlers` folder), `graphql` api saves us a lot of efforts
to write the Payload and Response and how parse them from the http `body`. `GraphQL` will handle them all for us automatically,
the `resolvers` and `Payload` and `Response` will be automatically generated using `gqlgen` library.

All the graphql resolvers and handlers can be located in the folder `graph` folder. For `session` and `user`, we will still
use the `REST` API. `graphql` will only serve the API for `todos` and `notes`, ofc, could be anything you want to serve
depending on your business logic requirements.

Run `make gqlgen` everytime you change the schema files (i.e. `*.graphql`) under the folder `./graph`. You can also go to
`gqlgen` website to find all the possible configurations and how to use it, see `./gqlgen.yml`. Note the `./graph/generated`
and `./graph/gmodel` are all generated code, do not change them. You should normally add the implementation into `./graph/*.resolvers.go`

You can get the `auth` token and then access `localhost:4000/playground` url to play with your graphql API.

# Test

Here we use `ginkgo` and `gomega` to manage our unit/integration test. For testing the http handlers and graphql API, we
use `httpexpect` to make it less verbose and more readable.

Run `make test` will run all the unit tests by `ginkgo`.

In the `test` mode, each test will spin up a new TestApp that forks a new `database` with a random database name and a new `port`.
After the whole test suite finishes, the database session pools will cleanup all the forked new database.
see `./db/sessionpool.go` to understand how it works. Each test will start from a fresh new database.

The benefit is that we can run all the tests independently in parallel without interfering with each other.

you could also run `ginkgo watch` to watch a folder and run `ginkgo -p` to run the tests in parallel to save your time.
You do not need to worry if the application could run in parallel, since every test will use their own database. Normally,
people will question and say "No, you are crazy, it is too slow". Before complaining the test performance, make a profiling
check first. It will take around 10 seconds if your run in parallel. If we are running the tests using docker in cloud,
performance is not at all a concern. We should care more about the test robustness and maintainability. A lot of mocking
the database will make it hard to evolve your database.

check `handlers, store, graph` folders for all the tests.

# Production & Staging

Staging and production normally run in the similar way, but just slightly different configs.
For real production release and configuration, it is very important that you do not leak your `JWT_SECRET`, otherwise
your website will be easily attacked. Ideally, we should not define those credentials inside the `helm` values either.
We should use something like Google's KMS (Key Management Service) to encrypt the config with your private key and GCP
will decrypt with your public key.

## Docker and docker-compose

Check the `dockerfile` on how to build the app into a docker and check the `.github/workflows/release.yml` on how to upload
your docker container to docker hub by `github action`.

We need to specify all the env variables correctly mentioned in config/config.go.
You could also check the docker-compose.yml to deploy using docker-compose if you want.

## Kubernetes

Kubernetes is the most popular way of running your app in production. Check `deploy` folder for deployment to kubernetes
using helm chart and helmfile. If you want to test and play with it locally, you could install `minikube`. Edit your
`/etc/hosts` to resolve the `domain` names in `ingres` with your IP address of `minikube`.

You could use `terraform` and `helmfile` to manage your infrastructure as code.
You could alternate the config for the production deployment and deploy to a new kubernetes for staging.

## Heroku

For a free staging deployment, I deploy to `heroku`'s free tier `app` and `postgres`. It is fairly easy to set it up.

- Create a new `app` in `heroku` following its official tutorial and then you can deploy your build docker to `heroku` app.
see the folder `.github/workflows/release.yml` where we use `github action` to deploy a docker to `heroku`. Note for gain
maximum flexibility, we do not use heroku's `go` app engine, instead we use the docker container stack.
- Manually create a `postgres` database in `heroku`.
- Config your `heroku`'s app's environment variables in the `settings` page. The `database` URL will automatically set
by `heroku`. You need to configure other environment variables.

