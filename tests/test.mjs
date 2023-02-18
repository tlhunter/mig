#!/usr/bin/env zx

// npm install -g zx
// forgive me for writing a Node.js-based integration test suite for a Golang project
// currently you need to destroy the migration tables before every test run

import assert from 'node:assert';

{
    console.log('### VERSION');

    const stdout = JSON.parse(await $`../../mig version --json`);

    assert.equal(typeof stdout.version, 'string', '.version is a string');
    assert.equal(typeof stdout.build_time, 'string', '.build_time is a string');
}

{
    console.log('### no subcommand');

    let didError = false;
    try {
        await $`../../mig --json`;
    } catch (out) {
        didError = true;
        assert.equal(out.exitCode, 1, 'exit status code 1');

        const stdout = JSON.parse(out.stdout);

        assert.equal(stdout.code, 'command_usage', 'right error code without arguments');
    }
    assert(didError, 'did in fact return an error');
}
    
{
    console.log('### STATUS (uninitialized)');

    let didError = false;
    try {
        await $`../../mig status --json --file="./test.migrc"`;
    } catch (out) {
        didError = true;
        assert.equal(out.exitCode, 9, 'exit status code 9');

        const stdout = JSON.parse(out.stdout);

        assert.equal(stdout.code, 'missing_tables', 'right error code without arguments');
    }

    assert(didError, 'did in fact return an error');
}

{
    console.log('### INIT');

    const stdout = JSON.parse(await $`../../mig init --json --file="./test.migrc"`);

    assert.equal(typeof stdout.success, 'string', '.success is a string');
    assert.ok(typeof stdout.success, '.success has a value');
}

{
    console.log('### STATUS (initialized)');
    const stdout = JSON.parse(await $`../../mig status --json --file="./test.migrc"`);

    assert.equal(stdout.locked, false, 'is not locked');
    assert.equal(typeof stdout.status, 'object', 'has .status');
    assert.equal(stdout.status.applied, 0, 'has no applied migrations yet');
    assert.equal(stdout.status.unapplied, 2, 'has two unapplied migrations');
    assert.equal(stdout.status.skipped, 0, 'has no skipped migrations');
    assert.equal(stdout.status.missing, 0, 'has no missing migrations');
    assert.equal(stdout.status.next, '20230101120058_add_users_table.sql', 'the first migration is the .next migration');
}

{
    console.log('### LIST');
    const stdout = JSON.parse(await $`../../mig list --json --file="./test.migrc"`);

    assert.ok(Array.isArray(stdout), 'output is an array');
    assert.equal(stdout.length, 2, 'output contains two entries');

    const addUsersTable = stdout[0];
    assert.equal(Object.keys(addUsersTable).length, 2, '[0] has two keys');
    assert.equal(Object.keys(addUsersTable.migration).length, 1, '[0].migration has one key');
    assert.equal(addUsersTable.migration.name, '20230101120058_add_users_table.sql', '[0].migration.name is correct');
    assert.equal(addUsersTable.status, 'unapplied', '[0].status is correct');

    const addEmailToUsers = stdout[1];
    assert.equal(Object.keys(addEmailToUsers).length, 2, '[0] has two keys');
    assert.equal(Object.keys(addEmailToUsers.migration).length, 1, '[0].migration has one key');
    assert.equal(addEmailToUsers.migration.name, '20230101120107_add_email_to_users.sql', '[0].migration.name is correct');
    assert.equal(addEmailToUsers.status, 'unapplied', '[0].status is correct');
}

{
    console.log('### LS');
    const list = await $`../../mig list --json --file="./test.migrc"`;
    const ls = await $`../../mig ls --json --file="./test.migrc"`;

    assert.equal(list.stdout, ls.stdout, 'list and ls are aliases');
}

{
    console.log('### LOCK (was unlocked)');
    const stdout = JSON.parse(await $`../../mig lock --json --file="./test.migrc"`);

    assert.equal(typeof stdout.success, 'string', '.success is a string');
    assert.ok(typeof stdout.success, '.success has a value');
    assert.ok(stdout.success.includes('success'), 'has a happy message');
}

{
    console.log('### STATUS');
    const stdout = JSON.parse(await $`../../mig status --json --file="./test.migrc"`);

    assert.equal(stdout.locked, true, 'is locked');
}

{
    console.log('### LOCK (was locked)');
    const stdout = JSON.parse(await $`../../mig lock --json --file="./test.migrc"`);

    assert.equal(typeof stdout.success, 'string', '.success is a string');
    assert.ok(typeof stdout.success, '.success has a value');
    assert.ok(stdout.success.includes('already'), 'has a mediocre message');
}

{
    console.log('### STATUS');
    const stdout = JSON.parse(await $`../../mig status --json --file="./test.migrc"`);

    assert.equal(stdout.locked, true, 'is locked');
}

{
    console.log('### UNLOCK (was locked)');
    const stdout = JSON.parse(await $`../../mig unlock --json --file="./test.migrc"`);

    assert.equal(typeof stdout.success, 'string', '.success is a string');
    assert.ok(typeof stdout.success, '.success has a value');
    assert.ok(stdout.success.includes('success'), 'has a happy message');
}

{
    console.log('### STATUS');
    const stdout = JSON.parse(await $`../../mig status --json --file="./test.migrc"`);

    assert.equal(stdout.locked, false, 'is not locked');
}

{
    console.log('### UNLOCK (was unlocked)');
    const stdout = JSON.parse(await $`../../mig unlock --json --file="./test.migrc"`);

    assert.equal(typeof stdout.success, 'string', '.success is a string');
    assert.ok(typeof stdout.success, '.success has a value');
    assert.ok(stdout.success.includes('already'), 'has a mediocre message');
}

{
    console.log('### STATUS');
    const stdout = JSON.parse(await $`../../mig status --json --file="./test.migrc"`);

    assert.equal(stdout.locked, false, 'is not locked');
}

{
    console.log('### UP');
    const stdout = JSON.parse(await $`../../mig up --json --file="./test.migrc"`);

    assert.equal(stdout.batch, 1, 'first ever batch');
    assert.equal(typeof stdout.migration, 'object', 'has .migration');

    assert.equal(stdout.migration.id, 1, 'first ever migration');
    assert.equal(stdout.migration.name, '20230101120058_add_users_table.sql', 'named after first migration file');
    assert.equal(stdout.migration.batch, 1, 'repeats the batch ID'); // silly
    assert.equal(typeof stdout.migration.time, 'string', 'has migration time');
}

{
    console.log('### LIST (one applied migration)');
    const stdout = JSON.parse(await $`../../mig list --json --file="./test.migrc"`);

    assert.ok(Array.isArray(stdout), 'output is an array');
    assert.equal(stdout.length, 2, 'output contains two entries');

    const addUsersTable = stdout[0];
    assert.equal(Object.keys(addUsersTable).length, 2, '[0] has two keys');
    assert.equal(Object.keys(addUsersTable.migration).length, 4, '[0].migration has four keys');
    assert.equal(addUsersTable.migration.id, 1, '[0].migration.id is correct');
    assert.equal(addUsersTable.migration.name, '20230101120058_add_users_table.sql', '[0].migration.name is correct');
    assert.equal(addUsersTable.migration.batch, 1, '[0].migration.batch is correct');
    assert.equal(typeof addUsersTable.migration.time, 'string', '[0].migration.time is present');
    assert.equal(addUsersTable.status, 'applied', '[0].status is correct');

    const addEmailToUsers = stdout[1];
    assert.equal(Object.keys(addEmailToUsers).length, 2, '[1] has two keys');
    assert.equal(Object.keys(addEmailToUsers.migration).length, 1, '[1].migration has one key');
    assert.equal(addEmailToUsers.migration.name, '20230101120107_add_email_to_users.sql', '[1].migration.name is correct');
    assert.equal(addEmailToUsers.status, 'unapplied', '[1].status is correct');
}

{
    console.log('### UP');
    const stdout = JSON.parse(await $`../../mig up --json --file="./test.migrc"`);

    assert.equal(stdout.batch, 2, 'second batch');
    assert.equal(typeof stdout.migration, 'object', 'has .migration');

    assert.equal(stdout.migration.id, 2, 'second migration');
    assert.equal(stdout.migration.name, '20230101120107_add_email_to_users.sql', 'named after second migration file');
    assert.equal(stdout.migration.batch, 2, 'repeats the batch ID'); // silly
    assert.equal(typeof stdout.migration.time, 'string', 'has migration time');
}

{
    console.log('### LIST (one applied migration)');
    const stdout = JSON.parse(await $`../../mig list --json --file="./test.migrc"`);

    assert.ok(Array.isArray(stdout), 'output is an array');
    assert.equal(stdout.length, 2, 'output contains two entries');

    const addEmailToUsers = stdout[1];
    assert.equal(Object.keys(addEmailToUsers).length, 2, '[1] has two keys');
    assert.equal(Object.keys(addEmailToUsers.migration).length, 4, '[1].migration has four keys');
    assert.equal(addEmailToUsers.migration.id, 2, '[1].migration.id is correct');
    assert.equal(addEmailToUsers.migration.name, '20230101120107_add_email_to_users.sql', '[0].migration.name is correct');
    assert.equal(addEmailToUsers.migration.batch, 2, '[1].migration.batch is correct');
    assert.equal(typeof addEmailToUsers.migration.time, 'string', '[1].migration.time is present');
    assert.equal(addEmailToUsers.status, 'applied', '[1].status is correct');
}
    
{
    console.log('### UP (none left)');

    let didError = false;
    try {
        await $`../../mig up --json --file="./test.migrc"`;
    } catch (out) {
        didError = true;
        assert.equal(out.exitCode, 1, 'exit status code 1');

        const stdout = JSON.parse(out.stdout);

        assert.equal(stdout.code, 'no_migrations', 'right error code without arguments');
    }

    assert(didError, 'did in fact return an error');
}

{
    console.log('### DOWN');
    const stdout = JSON.parse(await $`../../mig down --json --file="./test.migrc"`);

    assert.equal(typeof stdout.success, 'string');
    assert.ok(stdout.success.includes('20230101120107_add_email_to_users.sql'), 'mentions migration name');
}

{
    console.log('### LIST (one applied migration)');
    const stdout = JSON.parse(await $`../../mig list --json --file="./test.migrc"`);

    assert.ok(Array.isArray(stdout), 'output is an array');
    assert.equal(stdout.length, 2, 'output contains two entries');
    const addUsersTable = stdout[0];
    assert.equal(addUsersTable.status, 'applied', '[0].status is correct');

    const addEmailToUsers = stdout[1];
    assert.equal(Object.keys(addEmailToUsers).length, 2, '[1] has two keys');
    assert.equal(addEmailToUsers.migration.name, '20230101120107_add_email_to_users.sql', '[1].migration.name is correct');
    assert.equal(addEmailToUsers.status, 'unapplied', '[1].status is correct');
}

{
    console.log('### DOWN');
    const stdout = JSON.parse(await $`../../mig down --json --file="./test.migrc"`);

    assert.equal(typeof stdout.success, 'string');
    assert.ok(stdout.success.includes('20230101120058_add_users_table.sql'), 'mentions migration name');
}

{
    console.log('### LIST');
    const stdout = JSON.parse(await $`../../mig list --json --file="./test.migrc"`);

    assert.ok(Array.isArray(stdout), 'output is an array');
    assert.equal(stdout.length, 2, 'output contains two entries');

    const addUsersTable = stdout[0];
    assert.equal(Object.keys(addUsersTable).length, 2, '[0] has two keys');
    assert.equal(addUsersTable.status, 'unapplied', '[0].status is correct');

    const addEmailToUsers = stdout[1];
    assert.equal(addEmailToUsers.status, 'unapplied', '[1].status is correct');
}
    
{
    console.log('### DOWN (none left)');

    let didError = false;
    try {
        await $`../../mig down --json --file="./test.migrc"`;
    } catch (out) {
        didError = true;
        assert.equal(out.exitCode, 1, 'exit status code 1');

        const stdout = JSON.parse(out.stdout);

        assert.equal(stdout.code, 'nothing_to_revert', 'right error code without arguments');
    }

    assert(didError, 'did in fact return an error');
}

{
    console.log('### UPTO (fake)');

    let didError = false;
    try {
        await $`../../mig upto "2001_a_fake_migration.sql" --json --file="./test.migrc"`;
    } catch (out) {
        didError = true;
        assert.equal(out.exitCode, 1, 'exit status code 1');

        const stdout = JSON.parse(out.stdout);

        assert.equal(stdout.code, 'cannot_find_migration', 'right error code without arguments');
    }

    assert(didError, 'did in fact return an error');
}

{
    console.log('### UPTO (real)');
    const stdout = JSON.parse(await $`../../mig upto "20230101120107_add_email_to_users.sql" --json --file="./test.migrc"`);
        

    assert.ok(Array.isArray(stdout.migrations), '.migrations is an array');
    assert.equal(stdout.migrations.length, 2, '.migrations contains two entries');

    assert.equal(stdout.batch, 1, 'first batch since we started over');

    const addUsersTable = stdout.migrations[0];
    assert.equal(addUsersTable.id, 1, 'first migration is #1');
    assert.equal(addUsersTable.name, '20230101120058_add_users_table.sql', 'first migration file name');
    assert.equal(addUsersTable.batch, 1, 'first batch since we started over');

    const addEmailToUsers = stdout.migrations[1];
    assert.equal(addEmailToUsers.id, 2, 'second migration is #2');
    assert.equal(addEmailToUsers.name, '20230101120107_add_email_to_users.sql', 'second migration file name');
    assert.equal(addEmailToUsers.batch, 1, 'both migrations are in the first batch');
}

{
    console.log('### DOWN (reset to nothing)');

    await $`../../mig down --json --file="./test.migrc"`;
    await $`../../mig down --json --file="./test.migrc"`;
}

{
    console.log('### ALL');
    const stdout = JSON.parse(await $`../../mig all --json --file="./test.migrc"`);
        

    assert.ok(Array.isArray(stdout.migrations), '.migrations is an array');
    assert.equal(stdout.migrations.length, 2, '.migrations contains two entries');

    assert.equal(stdout.batch, 1, 'first batch since we started over');

    const addUsersTable = stdout.migrations[0];
    assert.equal(addUsersTable.id, 1, 'first migration is #1');
    assert.equal(addUsersTable.name, '20230101120058_add_users_table.sql', 'first migration file name');
    assert.equal(addUsersTable.batch, 1, 'first batch since we started over');

    const addEmailToUsers = stdout.migrations[1];
    assert.equal(addEmailToUsers.id, 2, 'second migration is #2');
    assert.equal(addEmailToUsers.name, '20230101120107_add_email_to_users.sql', 'second migration file name');
    assert.equal(addEmailToUsers.batch, 1, 'both migrations are in the first batch');
}

{
    console.log('### ALL (nothing to do)');

    let didError = false;
    try {
        await $`../../mig all --json --file="./test.migrc"`;
    } catch (out) {
        didError = true;
        assert.equal(out.exitCode, 1, 'exit status code 1');

        const stdout = JSON.parse(out.stdout);

        assert.equal(stdout.code, 'no_migrations', 'right error code without arguments');
    }
    assert(didError, 'did in fact return an error');
}