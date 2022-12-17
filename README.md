# mig

Platform agnostic binary for running database migrations.

Name to be determined.


## Backstory

I'm a Node.js developer. One package that I really like using is Knex, which is a Query Builder for Node.js. It's also a migration runner that exposes a nice pattern. It always bothered me that Knex was two things in one, namely a hard application dependency for the query builder aspect, while also being an optional dev dependency for executing migrations. It further bothered me when a project I was working on was tightly coupled to version X of the package for running thousands of migrations, and therefor making it exceedingly difficult to upgrade to version Y of the package for newer query builder capabilities.

I've worked on non-Node.js applications as well. Various projects seem to have their own migration runners. With Ruby and Python projects I would sometimes find myself writing Ruby or Python code in order to work with migration runners. Different projects end up having different migration runners and I would be forced to learn new patterns, many of them unoptimal. When things fail there's usually a stack trace for a language I don't care about. Running migrations requires that the right interpreter version is installed. Because of these reasons I've always wanted a generic, platform-agnostic, pre-compiled binary for running databaes migrations.

`mig` aims to be that tool.


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
