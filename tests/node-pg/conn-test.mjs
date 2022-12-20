import pg from 'pg';
const conn = {
  user: 'postgres',
  host: 'localhost',
  database: 'postgres',
  password: 'hunter2',
  port: 5432,
};
console.log(conn);
const client = new pg.Client(conn);

client.connect();

const res = await client.query('SELECT $1::text as message', ['db query successful']);
console.log(res.rows[0].message);
await client.end();
