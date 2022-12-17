# mig

Platform agnostic binary for running database migrations.

Name to be determined, usng `mig` as a placeholder.


## Backstory

I'm a Node.js developer. One package that I really like using is Knex, which is a Query Builder for Node.js. It's also a migration runner that exposes a nice pattern. It always bothered me that Knex was two things in one, namely a hard application dependency for the query builder aspect, while also being an optional dev dependency for executing migrations. It further bothered me when a project I was working on was tightly coupled to version X of the package for running thousands of migrations, and therefor making it exceedingly difficult to upgrade to version Y of the package for newer query builder capabilities.

I've worked on non-Node.js applications as well. Various projects seem to have their own migration runners. With Ruby and Python projects I would sometimes find myself writing Ruby or Python code in order to work with migration runners. Different projects end up having different migration runners and I would be forced to learn new patterns, many of them unoptimal. When things fail there's usually a stack trace for a language I don't care about. Running migrations requires that the right interpreter version is installed. Because of these reasons I've always wanted a generic, platform-agnostic, pre-compiled binary for running databaes migrations.

`mig` aims to be that tool.


## What is a Migration Runner?

A migration runner allows a developer to mutate a database schema in such a way that the mutations may be checked into a code repository. This is convenient because databaes mutations can be checked in alongside of code changes. It allows, say, SQL schema changes to be audited and visible and close to application code.

Not only can a database mutate forward, but it's also important to allow migrations to be undone. For this reason there are usually two separate sets or queries that get executed. One for the up/forward, and another for the down/back/rollback.


## Configuration

`mig` will be a single precompiled binary for running migrations on various platforms. Configuration can be achieved via environment variables, CLI flags, or even a config file. The config file will probably resemble a `.env` file with simple key=value pairs. `mig` will look in the current directory and traverse upwards until it finds a config file. Env vars will override the config file, and CLI flags will have the highest priority.

### Credentials

A SQL connection string is all we need for this. Basically it looks like `protocol://user:pass@host:port/dbname`.

```sh
mig --credentials="protocol://user:pass@host:port/dbname" init
MIG_CREDENTIALS="protocol://user:pass@host:port/dbname" mig init
```


## API

`mig` will support various flags and subcommands.

```sh
# list all migrations
mig ls

# create a migration named YYYY-MM-DD-HH-mm-ss-add_users_table.sql
mig create "Add users table"

# create the necessary migration tables
mig init

# run the next single migration, if it exists
mig runup

# run all of the unexecuted migrations, if any exist
mig runall

# run migrations up to and including the migration of this name
# if the named migration doesn't exist or isn't unexecuted then do nothing
mig runto YYYY-MM-DD-HH-mm-ss-add_users_table

# rolls back a single migration, prompting user to confirm, unless --force is provided
mig rundown --force

# rolls back migrations until the named migration is met, prompting user to confirm, unless --force is provided
mig rundownto YYYY-MM-DD-HH-mm-ss-add_users_table --force

# forcefully set / unset the lock, useful for fixing error scenarios
mig lock
mig unlock
```


## Tables

`mig` will probably require two tables. For this I'm taking inspiration from `Knex`. The two tables would be a table of migrations that have been executed, with a column for the name and a column for the time. The other table would be a simple locking mechanism to ensure multiple migrations don't run at once.


## Supported Platforms

Precompiled `mig` binaries will be provided for Linux, macOS, and Windows. At first they will be distributed by the releases feature of GitHub. These can, for example, be downloaded using `wget` or `curl` inside of a Dockerfile.

As far as DBMS go, I think it'll first support Postgres and will later add support for MySQL/MariaDB, SQLite, SQL Server, etc. A single binary will come with support for each of these databases to simplify things for the user.

Ideally the binary will be less than 10MB.


## Development

The language that `mig` is built in should only be of interest to those who want to contribute. Users of `mig` shouldn't care at all. One requirement is that `mig` is distributable as a compiled binary. I'm leaning towards Go, but certainly Rust would also be an option. Interpreted languages such as Node.js or Python or Ruby would be out of the question.


## Migration Files

Files can be created by hand or can be created with `mig create`. Files need to be uniquely named with an implicit order. The pattern that Knex uses is to prefix names with an ISO-8601 date and include a human-readable name for convenience. The timestamp is used for ordering and ensuring that migrations are executed in the proper order. In modern application development, users write code in parallel and check-in and merge then in non-deterministic order. For that reason an incrementing migration name doesn't work. Both engineers could run `i++` and get the same number.

Here are examples of migration filenames:

```
2022-12-11-12-15_create_users.sql
2022-12-14-12-15_create_projects.sql
2022-12-17-12-15_link_users_to_projects.sql
```

That said, `mig` also needs to support rollbacks of queries. If these files are purely SQL files then it probably requires that we have two separate files. For example, this could look like so:

```
2022-12-11-12-15_create_users.up.sql                 2022-12-11-12-15_create_users.down.sql
2022-12-14-12-15_create_projects.up.sql              2022-12-14-12-15_create_projects.down.sql
2022-12-17-12-15_link_users_to_projects.up.sql       2022-12-17-12-15_link_users_to_projects.down.sql
```

This is annoying because now two files need to be maintained. If a developer wants to change the name of one they need to change the name of the other. On the bright side SQL syntax highlighting still works perfectly for these files.

What Knex does is use JavaScript files that export an `up` and `down` function that get executed. This sucks because syntax highlighting doesn't work and the user must write langague-dependent code. On the bright side it's a single file solution.

Another approach would be to do something like create two SQL functions within each migration file, and these functions could then encompass the SQL queries that are to be run. This would then allow single file migrations for each up/down pair and maintain SQL highlighting. Most application developers don't know SQL function syntax but the `mig create` command handles creating files with the scaffolding present so it would be a non issue. This requires more research to define.
