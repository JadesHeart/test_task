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
