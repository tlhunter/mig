# `mig` the universal database migration runner

`mig` is a database migration runner that is distributed as precompiled binaries. The goal is to have a universal migration runner, one that is useful for projects written in any language. Gone are the days of learning a new technology when switching to a project written in Python or Node.js or Ruby. No longer sift through stack traces or install dependencies for languages that you don't usually work with. Simply download a binary and write SQL queries.

![mig list screenshot](./docs/screenshot-mig-list.png)

`mig` currently supports **PostgreSQL** and **MySQL** with plans to add more.


## v1.0 TODOs

`mig` isn't yet ready for production. When it is `mig` will hit version 1.0. Some of the required tasks to reach this include:

- [ ] implement `mig upto`
- [ ] address all of the TODOs
- [ ] unit test everything
- [ ] add support for sqlite
- [ ] support JSON output via `--json`
- [ ] write guides for migrating from other tools to `mig`
- [ ] make `mig create` templates for common scenarios like renaming columns
- [ ] allow developers to create their own template files


## What is a Migration Runner?

A migration runner allows a developer to mutate a database schema in such a way that the mutations may be checked into a code repository. This is convenient because database mutations can be reviewed and versioned alongside code changes. It allows, say, PostgreSQL schema changes to be audited and visible and close to application code.

Not only can a database mutate forward, but it's also important to allow migrations to be undone. For this reason there are usually two separate sets or queries that get executed. One for the up / forward operation, and another for the down / back / rollback / revert.


## Configuration

Configuration is be achieved by environment variables, CLI flags, or even a config file. The config file resembles a `.env` file with simple `KEY=value` pairs. While optional, `mig` looks in the current directory and traverses upwards until it finds a config file.

Configuration defined via CLI flags take the highest priority. After that are values in environment variables, and the configuration file has the lowest priority.

The variable names used in the `.migrc` file are named exactly the same as the environment variables. What follows is a list of the various configuration options.

### Credentials

A SQL connection string is all we need for this. Basically it looks like `protocol://user:pass@host:port/dbname`.

```sh
mig --credentials="protocol://user:pass@host:port/dbname"
MIG_CREDENTIALS="protocol://user:pass@host:port/dbname" mig
```

Currently, `mig` supports protocols of `postgresql` and `mysql`. In the future it will support more. Internally `mig` load the proper driver depending on the protocol. To disable TLS verification add the `?sslmode=disable` option. Here's an example of how you might connect to a local database:

```sh
mig --credentials="postgresql://user:hunter2@localhost:5432/dbname?sslmode=disable"
mig --credentials="mysql://user:hunter2@localhost:3306/dbname?tls=skip-verify"
```

### Migrations Directory

The migrations directory defaults to a folder named `migrations` in the current working directory. This can be overridden in the following ways:

```sh
mig --migrations="./db"
MIG_MIGRATIONS="./db" mig
```

### Configuration File Path

Unlike the other settings this one can only be set via CLI flag. To use it, specify a path to a `.migrc` file by using the `--flag` argument. This is useful for defining separate environments.

```sh
mig status --file="prod.migrc"
mig --file="local.migrc" down
```


## Commands

`mig` supports various commands:

```sh
# create the necessary migration tables
mig init

# get version information
mig version

# list all migrations
mig list

# check health of migrations, look for bugs, list unexecuted migrations
mig status

# create a migration named YYYYMMDDHHmmss_add_users_table.sql
mig create add_users_table
mig create "Add users table"

# run the next single migration
mig up

# TODO: run migrations up to and including the migration of this name
mig upto YYYYMMDDHHmmss_add_users_table

# run all of the unexecuted migrations
mig all

# rolls back a single migration
mig down

# forcefully set / unset the lock, useful for fixing error scenarios
mig lock
mig unlock
```


## Tables

`mig` requires two tables. This includes a table of migrations that have been executed and a simple locking mechanism ensuring multiple developers don't run migrations in parallel. These are created automatically by `mig init`.


## Migration File Syntax

Files are created by running `mig create`. Files need to be uniquely named and come with an implicit order. `mig` has chosen to use a number based on the time a migration was created. This is used to ensure that migrations are executed in the proper order. The filename is suffixed with a human-readable name for developer convenience.

Here are examples of migration filenames:

```
20221211093700_create_users.sql
20221214121500_create_projects.sql
20221217234100_link_users_to_projects.sql
```

We can think of changing a database schema as causing it to evolve. This evolution can be referred to as going "up". However, sometimes we'll create a migration that ends in disaster. When that happens we'll need to reverse this operation, referred to as going "down". For that reason a given migration file is made up of a pair of migrations: one up migration and one down migration.

> Note that generally a "down" migration is a destructive operation. Running them should only happen to recover from disaster. In fact, many teams that use database migrations only create "up" migrations.

These schema evolutions are represented as SQL queries. Often times they can be represented as a single query but in practice it's common to require multiple queries. Sometimes a given up migration just doesn't have a correlating down migration, so `mig` allows migrations to be empty. That means a migration is made up of zero or more queries.

To enable SQL syntax highlighting for migration files `mig` uses SQL comments to deliminate which queries are used for up the down migration.

Here's an example of a migration file:

```sql
--BEGIN MIGRATION UP--
CREATE TABLE user (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);
INSERT INTO user (id, name) VALUES (1, 'mig');
--END MIGRATION UP--

--BEGIN MIGRATION DOWN--
DELETE FROM user WHERE id = 1;
DROP TABLE user;
--END MIGRATION DOWN--
```

A migration file must contain one "up" migration block and one "down" migration block in that order. Any content outside of these two blocks is ignored. The queries that make up a block are executed in order and queries can span multiple lines. Be sure to end queries with a `;` semicolon.

Queries are wrapped in an implicit transaction since we don't want a migration to partially succeed. The transaction can be disabled by using a slightly different block syntax:

```sql
--BEGIN MIGRATION UP NO TRANSACTION--
CREATE TABLE accounts;
--END MIGRATION UP--

--BEGIN MIGRATION DOWN NO TRANSACTION--
DROP TABLE accounts;
--END MIGRATION DOWN--
```

Transactions should only be disabled when a situation calls for it, like when using `CREATE INDEX CONCURRENTLY`. When in doubt, leave transactions enabled.


## Development

Checkout the project then run the following commands to install dependencies, build, and run the program:

```sh
go get
make
./mig version
```

### Testing

The following commands let you easily spin up databases within Docker for testing:

```sh
docker run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=hunter2 -d postgres
docker run --name some-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=hunter2 -e MYSQL_DATABASE=mig -d mysql
```