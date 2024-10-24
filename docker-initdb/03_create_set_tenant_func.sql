CREATE OR REPLACE FUNCTION set_tenant(tenant_id INTEGER) 
RETURNS void AS $$
BEGIN
    PERFORM set_config('app.current_tenant', tenant_id::TEXT, false);
END;
$$ LANGUAGE plpgsql;