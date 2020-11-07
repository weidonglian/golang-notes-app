# Golang-notes-app

notes-app backend in golang. The app can run in three APP_MODE: [Dev, Test, Prod]. 
As default, it will run in `dev` mode. When you run test, it will switch to run in `test` mode.
`Prod` is a bit strict running environment, you need to define the required environment variables before running the main app.

See `./config/config.go` to see how to define those environments. 

It is easy to extend the APP_MODE and add your preferred config, e.g. [Staging].      

## Authentication and authorization

### Dual-tokens policy will be used here to mitigate the JWT's drawback of `logout`. 

- `access_token` will be used to access resources and has a short expired time duration. Once it is expired, the client 
side need to periodically fetch an access token using the `refresh_token` before the `access_token`'s expiration time. 
`access_token` will be stateless, we need to query the database to check the existence and validity.
- `refresh_token` will be used to access the `auth` end point to generate a new `access_token` and has a long expiration. 
The `refresh_token` will not be state-less. It needs to query the database if the `refresh_token` is still valid, 
when logout we need to set the `refresh_token` invalid or blacklisted.  

### JWT encryption algorithm.

Two choices: 

- `HS-256`. A simple one with only a secret key. Both token generation and verification use the same secret key, so-called
`symmetric`. If in your server side, you have control for both generation and verification on your own side, then it is 
your choice and should be secure and good enough. 
- `ES-256`. Use a key pair, private_key and public_key. The token generation will need your private key, and the verification
will use the public key to verify the signing and check the validity. 
 
# Development

## Requirements

In order to test develop locally, you need:

- Install `go`.
- Install `docker` for running postgres. 
- Add an `DEVROOT` in your `.bashrc` to define where to put the `dev` postgres volume data. See `Makefile` and search `DEVROOT`. 

After you have installed the above requirements. You are good to go.

## Useful commands

- `make start` will start the `main.go` app in dev mode.
- `make db-start-test` will start only the postgres docker.
- `make db-stop-test` will stop the postgres docker.
- `make db-ssh-test` will login the postgres docker container.

See the Makefile for more commands.

# Test

Run `make test` will run all the unit tests by `ginkgo`.

In the `test` mode, each test will spin up a new TestApp that forks a new `database` with a random database name and a new `port`.
After the whole test suite finishes, the database session pools will cleanup all the forked new database.
see `./db/sessionpool.go` to understand how it works. Each test will start from a fresh new database.

The benefit is that we can run all the tests independently in parallel without interfering with each other.    

# Production

We need to all the env variables correctly mentioned in config/config.go. 
You could also check the docker-compose.yml and `deploy` folder for deployment to kubernetes using helm chart and helmfile.

# Staging

You could alternate the config for the production deployment and deploy to a new kubernetes for staging. 
For a free staging deployment, I deploy to `heroku`'s free tier `app` and `postgres`. It is fairly easy to set it up.
   
- Create a new `app` in `heroku` following its official tutorial and then you can deploy your build docker to `heroku` app.
see the folder `.github/workflows/release.yml` where we use `github action` to deploy a docker to `heroku`. Note for gain
maximum flexibility, we do not use heroku's `go` app engine, instead we use the docker container stack. 
- Manually create a `postgres` database in `heroku`. 
- Config your `heroku`'s app's environment variables in the `settings` page. The `database` URL will automatically set
by `heroku`. You need to configure other environment variables.

For real production release and configuration, it is very important that you do not leak your `JWT_SECRET`, otherwise 
your website will be easily attacked.    