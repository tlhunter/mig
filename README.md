# `mig` the universal database migration runner

`mig` is a database migration runner that is distributed as precompiled binaries. The goal is to have a universal migration runner, one that is useful for projects written in any language. Gone are the days of learning a new technology when switching to a project written in Python or Node.js or Ruby. No longer sift through stack traces or install dependencies for unfamiliar languages. Simply download a binary and write SQL queries.

![mig list screenshot](./docs/screenshot-mig-list.png)

`mig` currently supports **PostgreSQL** and **MySQL** with plans to add more.


## v1.0 TODOs

`mig` isn't yet ready for production. When it is `mig` will hit version 1.0. Some of the required tasks to reach this include:

- [ ] document all exit status codes, ensure all errors return non-zero
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

Configuration is be achieved by environment variables, CLI flags, or even a config file. The config file resembles a `.env` file with simple `KEY=value` pairs. While optional, `mig` looks in the current working directory and traverses upwards until it finds a config file. A specific path can be specified using the `--file` flag.

Configuration defined via CLI flags take the highest priority. After that comes values in environment variables, with the configuration file having the lowest priority.

The variable names used in the `.migrc` configuration file use the same name as the environment variables. The following is a list of the various configuration options:

### Connection Credentials

A SQL connection string supplies all of the credentials. Basically it looks like `protocol://user:pass@host:port/dbname`.

```sh
mig --credentials="protocol://user:pass@host:port/dbname"
MIG_CREDENTIALS="protocol://user:pass@host:port/dbname" mig
```

Currently, `mig` supports protocols of `postgresql` and `mysql` with plans to support more. Internally `mig` loads the proper driver depending on the protocol. TLS checking can be set using query strings. Here's an example of how to connect to a local database:

```sh
mig --credentials="postgresql://user:hunter2@localhost:5432/dbname?tls=disable"
mig --credentials="mysql://user:hunter2@localhost:3306/dbname?tls=disable"
```

There are three connection string options for configuring secure databse connections:

* `?tls=verify`: enable secure database connection and verify the certificate
* `?tls=insecure`: enable secure database connection but don't verify (e.g. `localhost`)
* `?tls=disable` (default): use an insecure database connection

### Migrations Directory

The migrations directory defaults to `./migrations`. This can be overridden in the following ways:

```sh
mig --migrations="./db"
MIG_MIGRATIONS="./db" mig
```

### Configuration File Path

Unlike the other settings this one can only be set via CLI flag. To use it, specify a path to a config file by using the `--flag` argument. This is useful for defining separate environments. When specified, `mig` uses this path instead of searching for a `.migrc` file.

```sh
mig status --file="prod.migrc"
mig --file="local.migrc" down
```


## Commands

`mig` supports various commands:

| Command             | Purpose |
|---------------------|---------|
| `mig init`          | create the necessary migration tables |
| `mig version`       | display program version and compile time |
| `mig list`          | display a list of migrations, but finished and pending |
| `mig status`        | display health and status information |
| `mig create <name>` | creates a new migration file |
| `mig up`            | runs the next single migration |
| `mig upto <name>`   | runs migrations up to and including `<name>` |
| `mig all`           | runs all pending migrations |
| `mig down`          | rolls back the last executed migration |
| `mig unlock`        | unlocks the migration, in case of error |
| `mig lock`          | locks the migrations |

## Tables

`mig` requires two tables. This includes a table of migrations that have been executed and a simple locking mechanism ensuring multiple developers don't run migrations in parallel. These are created automatically by `mig init`.


## Migration File Syntax

Files are created by running `mig create`. Files need to be uniquely named and come with an implicit order. `mig` convention uses a number based on the time a migration was created. The filename is suffixed with a human-readable name for convenience.

Here are some example filenames:

```
20221211093700_create_users.sql
20221214121500_create_projects.sql
20221217234100_link_users_to_projects.sql
```

We can think of changing a database schema as causing it to evolve. This evolution can be referred to as going "up". However, sometimes we'll create a migration that ends in disaster. When that happens we'll need to reverse this operation, referred to as going "down". For that reason a given migration file is made up of a pair of migrations: one up migration and one down migration.

> Note that generally a "down" migration is a destructive operation. Running them should only happen to recover from disaster. In fact, many teams that use database migrations only create "up" migrations.

These schema evolutions are represented as SQL queries. Often times they can be represented as a single query but in practice it's common to require multiple queries. Sometimes a given up migration just doesn't conceptually have correlating down queries, so `mig` allows migrations to be empty. In other words a migration is made up of zero or more queries.

In order to enable SQL syntax highlighting for migration files `mig` uses SQL comments to deliminate which queries are used for the up and down migrations.

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

A migration file must contain one "up" migration block and one "down" migration block and in that order. Any content outside of these two blocks is ignored. The queries that make up a block are executed in order and queries can span multiple lines. Be sure to terminate queries with a `;` semicolon.

Queries are wrapped in an implicit transaction since we don't want a migration to partially succeed. The transaction can be disabled by using a slightly different block syntax:

```sql
--BEGIN MIGRATION UP NO TRANSACTION--
CREATE INDEX CONCURRENTLY foo_idx ON user (id, etc);
--END MIGRATION UP--

--BEGIN MIGRATION DOWN NO TRANSACTION--
DROP INDEX CONCURRENTLY foo_idx;
--END MIGRATION DOWN--
```

Transactions should only be disabled when a situation calls for it, like when using `CREATE INDEX CONCURRENTLY`. When in doubt, leave transactions enabled.


## Development

Clone the project then run the following commands to install dependencies, build a binary, and run the program:

```sh
go get
make
./mig version
```

### Testing

The following commands let you easily spin up databases within Docker for testing:

```sh
docker run --name some-postgres -p 5432:5432 \
  -e POSTGRES_PASSWORD=hunter2 -d postgres
docker run --name some-mysql -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=hunter2 -e MYSQL_DATABASE=mig -d mysql
```

The `tests/<DBMS>` directories contain config files and migrations to experiment with functionality.