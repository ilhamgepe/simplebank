-- simulasi deadlock jika account1 tf ke account 2 dan sebaliknya
--Tx1: transfer $10 from account 1 to account 2
BEGIN;

UPDATE
  accounts
SET
  balance = balance 10
WHERE
  id = 1 RETURNING *;

UPDATE
  accounts
SET
  balance = balance + 10
WHERE
  id = 2 RETURNING *;

ROLLBACK;


--Tx2: transfer $10 from account 2 to account 1
BEGIN;
UPDATE accounts SET balance = balance 10 WHERE id = 2 RETURNING *;
UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
ROLLBACK;
BEGIN;
INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES (1, 2, 10) RETURNING *;
INSERT INTO entries (account_id, amount) VALUES (1, -10) RETURNING *;
INSERT INTO entries (account_id, amount) VALUES (2, 10) RETURNING *;
SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
SELECT * FROM accounts WHERE id = 2 FOR UPDATE;
UPDATE accounts SET balance = =balance + 10 WHERE id = 2 RETURNING *;
ROLLBACK;



SELECT
  blocked_locks.pid AS blocked_pid,
  blocked_activity.usename AS blocked_user,
  blocking_locks.pid AS blocking_pid,
  blocking_activity.usename AS blocking_user,
  blocked_activity.query AS blocked_statement,
  blocking_activity.query AS current_statement_in_blocking_process,
  blocked_activity.application_name AS blocked_application,
  blocking_activity.application_name AS blocking_application
FROM
  pg_catalog.pg_locks blocked_locks
  JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
  JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
  AND blocking_locks.DATABASE IS NOT DISTINCT
FROM
  blocked_locks.DATABASE
  AND blocking_locks.relation IS NOT DISTINCT
FROM
  blocked_locks.relation
  AND blocking_locks.page IS NOT DISTINCT
FROM
  blocked_locks.page
  AND blocking_locks.tuple IS NOT DISTINCT
FROM
  blocked_locks.tuple
  AND blocking_locks.virtualxid IS NOT DISTINCT
FROM
  blocked_locks.virtualxid
  AND blocking_locks.transactionid IS NOT DISTINCT
FROM
  blocked_locks.transactionid
  AND blocking_locks.classid IS NOT DISTINCT
FROM
  blocked_locks.classid
  AND blocking_locks.objid IS NOT DISTINCT
FROM
  blocked_locks.objid
  AND blocking_locks.objsubid IS NOT DISTINCT
FROM
  blocked_locks.objsubid
  AND blocking_locks.pid != blocked_locks.pid
  JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE
  NOT blocked_locks.GRANTED;



SELECT
  a.datname,
  a.application_name,
  l.relation::regclass,
  l.transactionid,
  l.mode,
  l.locktype,
  l.GRANTED,
  a.usename,
  a.query,
  a.pid
FROM
  pg_stat_activity a
  JOIN pg_locks l ON l.pid = a.pid
where a.application_name = 'psql'
ORDER BY
  a.pid;
