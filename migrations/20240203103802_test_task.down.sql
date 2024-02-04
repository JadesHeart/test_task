BEGIN;

-- Удаление триггера
DROP TRIGGER IF EXISTS hash_password_trigger ON Users;

-- Удаление функции hash_password
DROP FUNCTION IF EXISTS hash_password();

-- Удаление расширения pgcrypto
DROP EXTENSION IF EXISTS pgcrypto;

-- Удаление таблицы Sessions
DROP TABLE IF EXISTS Sessions;

-- Удаление таблицы Users
DROP TABLE IF EXISTS Users;

COMMIT;
