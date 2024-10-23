/**
The owner of these tables should be 'root' as RLS do not apply to table owner. 
**/

-- Create the tenant table
CREATE TABLE app.tenant (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO app.tenant(id, name) VALUES (1, 'one'), (2, 'two');

-- Create the product table with tenant_id
CREATE TABLE app.product (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2)
);

INSERT INTO app.product (tenant_id, name, description, price) 
VALUES 
    (1, 'Product A1', 'Description for Product A1', 10.99),
    (1, 'Product A2', 'Description for Product A2', 15.50),
    (2, 'Product B1', 'Description for Product B1', 20.00),
    (2, 'Product B2', 'Description for Product B2', 25.99),
    (1, 'Product C1', 'Description for Product C1', 30.00),
    (1, 'Product C2', 'Description for Product C2', 45.00);

ALTER TABLE app.product ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation_policy ON app.product 
USING (tenant_id = current_setting('app.current_tenant')::INTEGER);

/**
Ensure the schema and tables are created before granting privileges.
**/
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA app TO dev;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA app TO dev;