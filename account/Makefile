.PHONY: migrate-create migrate-up migrate-down migrate-force create-keypair init

PWD = $(shell pwd)
MPATH = $(PWD)/migrations
PORT = 5432

# Default number of migrations to execute up or down
N = 1

# Commands for migrating tables
migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)

migrate-up:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable up $(N)

migrate-down:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable down $(N)

migrate-force:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable force $(VERSION)

# Commands to create private and public rsa 256 keys in account folder
create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(PWD)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(PWD)/rsa_private_$(ENV).pem -pubout -out $(PWD)/rsa_public_$(ENV).pem

# create dev and test keys
# run postgres containers in docker-compose
# migrate down
# migrate up
# docker-compose down
#init:
#	docker-compose up -d postgres-account && \
#	$(MAKE) create-keypair ENV=dev && \mod tidy
#
#	$(MAKE) create-keypair ENV=test && \
#	$(MAKE) migrate-down APPPATH=account N= && \
#	$(MAKE) migrate-up APPPATH=account N= && \
#	docker-compose down