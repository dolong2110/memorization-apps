# Memorization App
[![Actions Status](https://github.com/dolong2110/Memoirization-Apps/workflows/build/badge.svg)](https://github.com/dolong2110/Memoirization-Apps/actions)
[![codecov](https://codecov.io/gh/dolong2110/Memoirization-Apps/branch/master/graph/badge.svg)](https://codecov.io/gh/dolong2110/Memoirization-Apps)


## Application Overview

A chart of the tools and applications used is given below.

![App Overview](pictures/application_overview.png)

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