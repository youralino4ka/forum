CREATE OR REPLACE FUNCTION delete_old_messages()
RETURNS VOID AS $$
BEGIN
    DELETE FROM messages WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;