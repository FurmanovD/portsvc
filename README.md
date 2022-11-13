# Port service

## Business problem
Ports service is aimed to create&update port records in the database.

## Config
Service's configuration example is in the file: **./config/example-config.ini**

**NOTE:** It contains configuration **secrets** that **must not be used** in production and should be changed to loading secrets from external secure storage like Hashicorp Vault.

To run the service locally, create a copy of **example-config.ini** and rename it to **config.ini** or use a command-line parameter `-config`. e.g. `-config=./config/config-local.ini`

## Build
`make build` build a **builder** docker container, builds a binary file in it and constructs a service docker image basing on alpine linux image.

## Test
Run tests locally: `make test`
The only unit-tests present are in `pkg/commoncfg/mysql_test.go` to demonstrate usage of `testing/quick` package.

Run tests in a docker container to avoid local environment peculiarities: `make test-docker`

## Run
To spin-up a docker compose, make sure **config.ini** is identical to **example-config.ini** and run: `make run`
It will spin up the docker containers of a Database (MySQL) and a Port service.

## Ingest ports file
To make service ingest the ports file, copy it into a directory `/in` of the running `portsvc` container: e.g. `docker cp /tmp/ports.txt portsvc:/in`

## DB
To access the database, you can use address: `localhost:3306`

## Initial DB data
Run DB container: `make run-deps`

Using any database/sql client, apply scripts from files: **internal/pkg/db/migrations/1_schema.sql**
MySQL container will bind the local **./test/db/** folder to save all its data, so you can stop the system at any time - the data will persist.

To debug from IDE, you can run DB container alone by run command: `make run-deps`

## Stop
To stop all the fun, run: `make stop`
