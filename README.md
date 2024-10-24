# multi-tenancy-with-RLS
Demonstration of Multi-Tenant SaaS application using Row Level Security in Postgres. 

This is an example repository of achieving Multi-Tenancy in Golang using Row Level Security (RLS) for article [rls](eddie023.github.io/post/rls). 

# Overview 
1. All the necessary database SQL queries such as creating tables, inserting seeds, and creating a 'dev' role are done via database initialization script.
2. Use config.beforeAcquire and config.AfterRelease to set/reset tenant_id configuration parameter in the DB. 

# Running the app 
1. Run `docker compose up`
