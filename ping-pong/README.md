📘 Database Setup in Kubernetes (Manual SQL Initialization)

This document explains how to manually initialize the PostgreSQL database running inside the Kubernetes cluster.
This method is useful during development or when you need to quickly bootstrap your tables without using migration tools or automation.

Note:
This is a quick and manual method. For production environments, automated migrations (e.g., Kubernetes Jobs or a migration CI pipeline) are recommended.

🚀 Overview

We have deployed PostgreSQL into the Kubernetes cluster using a StatefulSet:
```
Namespace: exercises

StatefulSet: postgres-stset

Service: postgres-svc

Database name: pingpong

User: stac

Password: password
```

PostgreSQL automatically creates the user and database on first startup based on the environment variables:
```
env:
  - name: POSTGRES_USER
    value: "stac"
  - name: POSTGRES_PASSWORD
    value: "password"
  - name: POSTGRES_DB
    value: "pingpong"
```

However, tables and schema are NOT created automatically.
We must create them manually using SQL.

🧩 Step 1 — Exec Into the PostgreSQL Pod

Find the pod:

`kubectl get pods -n exercises`


It should look like:

`postgres-stset-0`


Open an interactive shell inside the database pod:

`kubectl exec -it -n exercises postgres-stset-0 -- sh`

🧩 Step 2 — Connect to PostgreSQL

Once inside the pod, connect to the pingpong database using the default user:

`psql -U stac pingpong`


If the connection succeeds, you will see:

`pingpong=#`

🧩 Step 3 — Create Required Tables

Run the SQL statements directly inside the psql console.

Create the lines table:
```
CREATE TABLE IF NOT EXISTS lines (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

Confirm the table exists:

`\d lines`


You should see the schema.

🧪 Step 4 — Test the Application

You can now run your application.
The table exists, so queries like INSERT or SELECT should work without errors.

Check logs:

kubectl logs -n exercises <your-app-pod>


You should no longer see:

pq: relation "lines" does not exist

🧹 Optional — Reset the Database

If you want to drop the table during development:

`DROP TABLE lines;`


Or drop and recreate the whole database (inside postgres connection):
```
psql -U stac
DROP DATABASE pingpong;
CREATE DATABASE pingpong;
```

`k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2`