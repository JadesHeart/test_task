BEGIN;

CREATE TABLE IF NOT EXISTS Users (
                                     UserID bigserial not null PRIMARY KEY,
                                     Username VARCHAR(255) NOT NULL UNIQUE,
    PasswordHash VARCHAR(255) NOT NULL,
    Salt VARCHAR(255) NOT NULL,
    FailedLoginAttempts INT NOT NULL DEFAULT 0
    );


CREATE TABLE IF NOT EXISTS Sessions(
                                       SessionID bigserial PRIMARY KEY,
                                       UserID INT REFERENCES Users(UserID),
    Token VARCHAR(255) NOT NULL,
    ExpirationTime TIMESTAMPTZ NOT NULL
    );


CREATE EXTENSION IF NOT EXISTS pgcrypto;



CREATE OR REPLACE FUNCTION generate_salt(length INT) RETURNS VARCHAR AS $$
BEGIN
RETURN substring(md5(random()::text || clock_timestamp()::text), 1, length);
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION hash_password() RETURNS TRIGGER AS $$
BEGIN
  NEW.salt := generate_salt(16);
  NEW.passwordhash := md5(NEW.PasswordHash || NEW.salt);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;



CREATE TRIGGER hash_password_trigger
    BEFORE INSERT ON Users
    FOR EACH ROW
    EXECUTE FUNCTION hash_password();


COMMIT;
