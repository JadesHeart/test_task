CREATE TRIGGER hash_password_trigger
    BEFORE INSERT ON Users
    FOR EACH ROW
    EXECUTE FUNCTION hash_password();
