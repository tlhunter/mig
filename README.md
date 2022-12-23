# `mig`

Platform agnostic binary for running database migrations.


## Backstory

I'm predominantly a Node.js developer. One package that I really like using is Knex, which is a Query Builder for Node.js. It's also a migration runner that exposes a nice pattern. It always bothered me that Knex was two things in one, namely a hard application dependency for the query builder aspect, while also being an optional dev dependency for executing migrations. It further bothered me when a project I was working on was tightly coupled to version X of the package for running thousands of migrations, and therefor making it exceedingly difficult to upgrade to version Y of the package for newer query builder capabilities.

I've worked on non-Node.js applications as well. Various projects seem to have their own migration runners. With Ruby and Python projects I would sometimes find myself writing Ruby or Python code in order to work with migration runners. Different projects end up having different migration runners and I would be forced to learn new patterns, many of them unoptimal. When things fail there's usually a stack trace for a language I don't care about. Running migrations requires that the right interpreter version is installed. Because of these reasons I've always wanted a generic, platform-agnostic, pre-compiled binary for running databaes migrations.

`mig` aims to be that tool.


## v0.1 Progress

- [X] parse env vars
- [X] parse CLI flags
- [X] parse config files
- [X] connect to PG database
- [X] decide on a migration file format
- [X] implement `mig create`
- [X] implement `mig init`
- [X] implement `mig lock` and `mig unlock`
- [ ] implement `mig list`
- [ ] implement `mig status`
- [ ] implement `mig runup` and `mig rundown`
- [ ] implement `mig runto` and `mig rundownto`

## v1.0 Progress

- [ ] address all of the TODOs
- [ ] unit test everything
- [ ] syntax for disabling migration transactions
- [ ] allow specifying path to config file via `--file`
- [ ] add support for mysql
- [ ] add support for sqlite
- [ ] write guides for migrating from other tools to `mig`


## What is a Migration Runner?

A migration runner allows a developer to mutate a database schema in such a way that the mutations may be checked into a code repository. This is convenient because databaes mutations can be checked in alongside of code changes. It allows, say, SQL schema changes to be audited and visible and close to application code.

Not only can a database mutate forward, but it's also important to allow migrations to be undone. For this reason there are usually two separate sets or queries that get executed. One for the up/forward, and another for the down/back/rollback.


## Configuration

`mig` will be a single precompiled binary for running migrations on various platforms. Configuration can be achieved via environment variables, CLI flags, or even a config file. The config file resembles a `.env` file with simple `key=value` pairs. `mig` will look in the current directory and traverse upwards until it finds a config file. Configuration priority follows this order:

- CLI Flags
- Env Vars
- Config File

Configuration contains at least the following data:

* Connection string
* Migrations directory

### Credentials

A SQL connection string is all we need for this. Basically it looks like `protocol://user:pass@host:port/dbname`.

```sh
mig --credentials="protocol://user:pass@host:port/dbname" init
MIG_CREDENTIALS="protocol://user:pass@host:port/dbname" mig init
```


## API

`mig` will support various flags and commands.

```sh
# create the necessary migration tables
mig init

# list all migrations
mig list # or mig ls

# check health of migrations, look for bugs, list unexecuted migrations
mig status

# create a migration named YYYYMMDDHHmmss-add_users_table.sql
mig create "Add users table"

# run the next single migration, if it exists
mig runup

# run all of the unexecuted migrations, if any exist
mig runall

# run migrations up to and including the migration of this name
# if the named migration doesn't exist or isn't unexecuted then do nothing
mig runto YYYYMMDDHHmmss-add_users_table

# rolls back a single migration, prompting user to confirm, unless --force is provided
mig rundown --force

# rolls back migrations until the named migration is met, prompting user to confirm, unless --force is provided
mig rundownto YYYYMMDDHHmmss-add_users_table --force

# forcefully set / unset the lock, useful for fixing error scenarios
mig lock
mig unlock
```

Common flags include:

```sh
mig --connection 'connection string'
mig --migrations './migrations'
mig --file prod.migrc
mig --debug
```


## Tables

`mig` requires two tables, taking inspiration from `Knex`. This includes a table of migrations that have been executed. The other table would be a simple locking mechanism to ensure multiple migrations don't run at once.


## Supported Platforms

Precompiled `mig` binaries will be provided for Linux, macOS, and Windows. At first they will be distributed by the releases feature of GitHub. These can, for example, be downloaded using `wget` or `curl` inside of a Dockerfile.

As far as DBMS go, I think it'll first support Postgres and will later add support for MySQL/MariaDB, SQLite, SQL Server, etc. A single binary will come with support for each of these databases to simplify things for the user.

Ideally the binary will be less than 10MB.


## Migration Files

Files can be created by hand or can be created with `mig create`. Files need to be uniquely named with an implicit order. `mig` has chosen to use a number based on the ISO-8601 date/time standard. This is used to ensure that migrations are executed in the proper order. This time is suffixed with a human-readable name for developer convenience.

In modern application development developers write code in parallel and check-in and merge them in non-deterministic order. For that reason an incrementing integer migration name just doesn't work. For example if migration "17" is checked-in and two engineers increment it they'll both end up with "18".

Here are examples of migration filenames:

```
20221211093700-create_users.sql
20221214121500-create_projects.sql
20221217234100-link_users_to_projects.sql
```

When modifying a database schema we can think of it as evolving the database. This evolution can be referred to as going "up". However, sometimes we'll create a migration that ends in disaster. When that happens we'll need to reverse this operation, referred to as going "down". For that reason a given migration file is made up of a pair of migrations: one up migration and one down migration.

_Note that generally a "down" migration is a destructive operation. Running them should only happen to recover from disaster. In fact, many teams that use databaes migrations only crate "up" migrations._

These schema evolutions are represented as SQL queries. Often times they can be represented as a single query but in practice it's very common to require multiple queries. Sometimes a given up migration just can't be represented with a down migration, so we allow a migration to be empty as well. For this reason we say that a migration is made up of zero or more queries.

In order to allow SQL syntax highlighting to play nicely with migration files we'll make use of SQL comments to deliminate which part of the files is the up or the down migration.

Here's an example of a migration file:

```sql
--BEGIN MIGRATION UP--
CREATE TABLE user (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);
--END MIGRATION UP--

--BEGIN MIGRATION DOWN--
DROP TABLE user;
--END MIGRATION DOWN--
```

A migration file must contain one up migration block and one down migration block, and in that order. Any content outside of these two blocks is ignored. The queries that make up a block are executed in order and queries can span multiple lines. Queries are wrapped in an implicit transaction since we don't want a migration to only be executed partially.