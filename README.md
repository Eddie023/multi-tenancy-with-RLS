# multi-tenancy-with-RLS
This is an example repository of achieving Multi-Tenancy in Golang using Row Level Security (RLS) in Postgres. Companion article: [rls](eddie023.github.io/post/rls). 

# Overview 
1. `docker-initdb` contains all the necessary SQL commands required to initialize our containarized database.
2. For demonstration purposes, `tenantId` is dynamically set via query param and updated via request context. 
3. While configuring the database connection, we are dynamically setting the `app.current_tenant` configuration paramter using `set_tenant` SQL function using `BeforeAcquire` hook from pgx library. 
4. Similarily, we will reset the tenantId using `AfterRelease` hook. 

# Running the application
1. Run `docker compose up` 
2. Wait few seconds for `apiserver` to be live as it waits until database is intialized. 
3. Run `curl -X GET http://localhost:8848/products\?tenantId=<1|2>` to dynamically set the tenantId values.  
