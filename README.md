# Memorization App
[![Actions Status](https://github.com/dolong2110/memorization-apps/workflows/Build%20and%20Test/badge.svg)](https://github.com/dolong2110/memorization-apps/actions)
[![CI](https://github.com/dolong2110/memorization-apps/workflows/Continuous-Integration/badge.svg)](https://github.com/dolong2110/memorization-apps/actions?query=workflow%3CI)
[![codecov](https://codecov.io/gh/dolong2110/memorization-apps/branch/master/graph/badge.svg?token=2a2ab5db-5712-4668-9753-c55d541550fb)](https://codecov.io/gh/dolong2110/memorization-apps)
[![Go Report Card](https://goreportcard.com/badge/github.com/dolong2110/memorization-apps/account)](https://goreportcard.com/report/github.com/dolong2110/memorization-apps/account)

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