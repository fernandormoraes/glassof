# Glassof 🍷

Glassof is a CLI application for creating a replication data service.

## About

Main porpuse of Glassof is to serve data from PostgreSQL database to MongoDB database in a easily way.

Take a ./glassof init to prepare PostgreSQL Database for replication

## Development state
Glassof is in his early initial development phase, this means that features are not ready even for testing porpuses.

## Dependencies
Actually it uses **wal2json** to get logical replication output for replication to MongoDB, so it'll not work until the extension is installed.

## Concept

The concept of glassof is to setup PostgreSQL Database for logical replication and configure a service to replicate data to MongoDB in a easily way using CLI commands with features as auto-create services for replication data between PostgreSQL and MongoDB, preparing PostgreSQL for logical replication using wal2json extension, listening for changes and CRUD everything into a MongoDB collection with custom data transformation aswell.


*So seat and take a glassof replicated data.*
