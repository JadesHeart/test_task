BEGIN;

CREATE TABLE IF NOT EXISTS Users (
                                     UserID bigserial not null PRIMARY KEY,
                                     Username VARCHAR(255) NOT NULL UNIQUE,
    PasswordHash VARCHAR(255) NOT NULL,
    Salt VARCHAR(255) NOT NULL,
    FailedLoginAttempts INT NOT NULL DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS Sessions (
                                        SessionID bigserial PRIMARY KEY,
                                        UserID INT REFERENCES Users(UserID),
    Token VARCHAR(255) NOT NULL,
    ExpirationTime TIMESTAMPTZ NOT NULL
    );

COMMIT;