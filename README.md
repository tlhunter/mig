# `mig` the universal database migration runner

`mig` is a database migration runner distributed as precompiled binaries. The goal is to have a universal migration runner that is useful for projects written in any language. Gone are the days of learning a new migration tool when switching to a project written in Python or Node.js or Ruby. No longer sift through stack traces or install dependencies for unfamiliar languages. Simply download a binary and write SQL queries.

![mig list screenshot](./docs/screenshot-mig-list.png)

`mig` currently supports **PostgreSQL**, **MySQL**, and **SQLite**. Pull requests to support other databases are appreciated.

`mig` will be released as version 1.0 when it is ready for production use.


## What is a Migration Runner?

A migration runner allows a developer to mutate a database schema in such a way that the mutations may be checked into a code repository. This conveniently allows database mutations to be reviewed and versioned alongside relevant code changes. Both schema and code audits can happen within a single pull request.

Not only can a database mutate forward, but it's also important to allow migrations to be undone. For this reason there are usually two separate sets or queries as part of a given migration. One for the up / forward operation, and another for the down / back / rollback / revert operation.


## Configuration

Configuration can be achieved by environment variables, CLI flags, config file, or combinations of these. The config file resembles a `.env` file, containing simple `KEY=value` pairs. While optional, `mig` looks in the current working directory and traverses upwards until it finds a config file. Alternatively a path can be provided with the `--file` flag.

Configuration defined via CLI flags take the highest priority. After that comes values in environment variables, with the configuration file having the lowest priority.

The variable names used in the `.migrc` configuration file use the same name as the environment variables. The following is a list of configuration options:

### Connection Credentials

A SQL connection string supplies all of the credentials. It resembles a URL format of `protocol://user:pass@host:port/dbname`:

```sh
mig --credentials="protocol://user:pass@host:port/dbname"
MIG_CREDENTIALS="protocol://user:pass@host:port/dbname" mig
```

`mig` currently supports protocols of `postgresql`, `mysql`, and `sqlite`. Internally `mig` loads the proper driver depending on the protocol. TLS checking can be set using query strings. Here's an example of how to connect to a local database for these different database systems:

```sh
mig --credentials="postgresql://user:hunter2@localhost:5432/dbname?tls=disable"
mig --credentials="mysql://user:hunter2@localhost:3306/dbname?tls=disable"
mig --credentials="sqlite://localhost/path/to/database.db"
```

> Note that SQLite follows a slightly different format. Paths have the first slash stripped. This means `localhost/foo.db` translates to `foo.db`, `localhost//tmp/foo.db` (with two slashes) becomes `/tmp/foo.db`, `localhost/../foo.db` becomes `../foo.db`, etc. This format is awkward and will likely change.

There are three connection string options for configuring secure database connections:

* `?tls=verify`: enable secure database connection and verify the certificate
* `?tls=insecure`: enable secure database connection but don't verify (e.g. `localhost`)
* `?tls=disable` (default): use an insecure database connection

These values are appended to the end of the connection URL.

### Migrations Directory

The migrations directory defaults to `./migrations` but can be overridden:

```sh
mig --migrations="./db"
MIG_MIGRATIONS="./db" mig
```

### Configuration File Path

Unlike the other settings this one can only be set via CLI flag. To use it, specify a path to a config file with the `--file` flag. This is useful for declaring separate environments. When specified, `mig` uses this path instead of searching for a `.migrc` file:

```sh
mig status --file="prod.migrc"
mig --file="local.migrc" down
```

### JSON Output

Provide the `--json` flag and `mig` will output a single line valid JSON document:

```sh
mig version --json
```
```json
{"version":"0.4.0","build_time":"2023-02-20T17:23:12-0800"}
```

This setting is only supported by CLI flag and has no environment variable or config file equivalent.


## Commands

`mig` supports various commands:

| Command             | Purpose |
|---------------------|---------|
| `mig init`          | create the necessary migration tables |
| `mig version`       | display program version and compile time |
| `mig list`          | display a list of migrations, including finished and pending |
| `mig status`        | display health and status information |
| `mig create <name>` | create a new migration file named `<name>` |
| `mig up`            | run the next single migration |
| `mig upto <name>`   | run migrations up to and including `<name>` |
| `mig all`           | run all pending migrations |
| `mig down`          | rolls back the last executed migration |
| `mig unlock`        | unlock migrations, usually to recover from an error |
| `mig lock`          | manually lock migrations |

## Tables

`mig` requires two tables. This includes a table of migrations that have been executed and a simple locking mechanism ensuring multiple developers don't run migrations in parallel. These are created automatically by `mig init`.

`mig` automatically sets the lock when performing operations that modify the database. This includes the `up`, `upto`, `all`, and `down` commands. Once the command completes successfully it is unlocked. If a migration fails then the lock will remain so that an engineer knows to investigate.


## Migration File Syntax

Migration files are created by running `mig create`. Files need to be uniquely named and come with an implicit order. `mig` convention uses a number prefix based on the time a migration was created to guarantee uniqueness and order. Filenames are suffixed with a human-readable title for convenience.

Here are some example filenames:

```
20221211093700_create_users.sql
20221214121500_create_projects.sql
20221217234100_link_users_to_projects.sql
```

We can think of changing a database schema as causing it to evolve. This evolution can be referred to as going "up". However, sometimes we'll create a migration that ends in disaster. When that happens we'll need to reverse this operation, referred to as going "down". For that reason a given migration file is made up of a pair of migrations: one up migration and one down migration.

> Note that generally a "down" migration is a destructive operation. Running them should only happen to recover from disaster. In fact, many teams that use database migrations only create "up" migrations, opting to perform manual rollbacks.

These schema evolutions are represented as SQL queries. Often times they can be represented as a single query but in practice it's common to require multiple queries. Sometimes a given up migration just doesn't conceptually have correlating down queries, so `mig` allows migrations to be empty. In other words a migration is made up of zero or more queries.

Migration files are written in pure SQL and use SQL comments for deliminating which queries belong to the up or down migration. This was done to enable syntax highlighting.

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

By default queries are wrapped in an implicit transaction since we don't want a migration to partially succeed. The transaction can be disabled by using a slightly different syntax:

```sql
--BEGIN MIGRATION UP NO TRANSACTION--
CREATE INDEX CONCURRENTLY foo_idx ON user (id, etc);
--END MIGRATION UP--

--BEGIN MIGRATION DOWN NO TRANSACTION--
DROP INDEX CONCURRENTLY foo_idx;
--END MIGRATION DOWN--
```

Transactions should only be disabled when a situation calls for it, like when using `CREATE INDEX CONCURRENTLY`. When in doubt, leave transactions enabled. Consider breaking up a complex migrations that contain queries that should run with and without a transaction.


## Contributing

Clone the project then run the following commands to install dependencies, build a binary, and run the program:

```sh
go get
make
./mig version
```

### Testing

There are some unit tests. Those can be run with the following command:

```sh
make test
```

However most of the tests are available as integration tests which require active databases.

The following commands spin up databases within Docker for local testing:

```sh
docker run --name some-postgres -p 5432:5432 \
  -e POSTGRES_PASSWORD=hunter2 -d postgres
docker run --name some-mysql -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=hunter2 -e MYSQL_DATABASE=mig -d mysql
```

The `tests/<database>` directories contain config files and migrations to experiment with functionality.

Run the following commands to run the different integration test suites:

```sh
npm install -g zx

cd tests/postgres # or mysql or sqlite
node ../test.mjs
```

Unfortunately the integration tests currently require Node.js to be installed. This will be fixed in the future.

For a fast local development cycle, go into the `tests/sqlite` directory and then execute the following one liner after you've made code changes:

```sh
clear ; pushd ../.. && make && popd && rm test.db ; ../test.mjs
```

This recompiles the binary, deletes the sqlite database file (thus resetting the state of the database), and runs the suite of acceptance tests.