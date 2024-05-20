# Utils

It's under development.

## Description

This project is a set of functionalities such as:

- A database type to open a connection and manage database migrations easily.
- A database repository on top of [GORM](https://gorm.io/) that provides easier transaction management as well as common methods like `Save` or `Find`.
- An event bus type to connnect to a [NATS](https://nats.io/) message queue.
- An utility to load `.env` files.
- An HTTP library to run HTTP requests.

## Usage

- [Database orm](./udatabase/uorm/README.md).
- [Database sql](./udatabase/usql/README.md).
- [HTTP requester](./requester/README.md).
