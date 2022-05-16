# Memorization App


## Application Overview

A chart of the tools and applications used is given below.

![App Overview](./application_overview.png)

## Use Makefile to generate rsa256 keys

````
make create-keypair ENV=test
````

## Using Docker

each time change the code we should re-initialize the docker again at the root.

````
docker-compose up
````

can run a container that does not depend on other

````
docker-compose up postgres-account
````

## Migrate DB

````
make migrate-create NAME=add_users_table
make migrate-up // update table
make migrate-down // revert table
````