.PHONY: generate build

SQL_SOURCES := $(wildcard sql/queries/*.sql) $(wildcard sql/schema/*.sql) sqlc.yaml

generate:	$(SQL_SOURCES)
	go tool sqlc generate

build:	generate 
	go build

