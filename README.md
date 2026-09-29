# Cat-Dog Application

A containerized Cat/Dog web application consisting of two Dockerized Go APIs, PostgreSQL, and an Nginx reverse proxy.

## Architecture

The application is split across two Ubuntu Server VMs.

```text
                         VM1
              ┌─────────────────────┐
              │ Dockerized Nginx    │
              │ :80 / :443          │
              │                     │
              │ cats.example.local  │
              │ dogs.example.local  │
              │                     │
              │ Static frontend     │
              └──────────┬──────────┘
                         │
                         │ HTTP
                         ▼
                         VM2
              ┌─────────────────────┐
              │ Docker Compose      │
              │                     │
              │ cat-api :8080       │
              │ dog-api :8081       │
              │ postgres :5432       │
              └─────────────────────┘
```

### VM1

Runs the Nginx container and serves the frontend websites.

* Cat site: `cats.example.local`
* Dog site: `dogs.example.local`
* HTTP: `80`
* HTTPS: `443`

Nginx also acts as a reverse proxy:

```text
/api/cat → VM2:8080/cat
/api/dog → VM2:8081/dog
```

### VM2

Runs the application backend using Docker Compose:

* Cat API: `8080`
* Dog API: `8081`
* PostgreSQL: `5432` inside the Docker network
* PostgreSQL host port: `5433`

The APIs connect to PostgreSQL using the Docker Compose service name:

```text
postgres:5432
```

## Repository Structure

```text
cat-dog-app/
├── README.md
├── .gitignore
│
├── vm1-nginx/
│   ├── docker-compose.yml
│   ├── nginx.conf
│   ├── conf.d/
│   │   └── cats-dogs.conf
│   ├── ssl/
│   │   └── cats-dogs.crt
│   └── www/
│       ├── cat-site/
│       │   ├── index.html
│       │   ├── app.js
│       │   └── style.css
│       └── dog-site/
│           ├── index.html
│           ├── app.js
│           └── style.css
│
└── vm2-app/
    ├── docker-compose.yml
    ├── .env.example
    ├── cat-api/
    │   ├── Dockerfile
    │   ├── go.mod
    │   ├── go.sum
    │   └── main.go
    ├── dog-api/
    │   ├── Dockerfile
    │   ├── go.mod
    │   ├── go.sum
    │   └── main.go
    └── postgres/
        └── init/
            ├── 01-cats-schema.sql
            ├── 02-cats-seed.sql
            ├── 03-dogs-schema.sql
            └── 04-dogs-seed.sql
```

## Prerequisites

Both VMs should have:

* Ubuntu Server
* Docker
* Docker Compose plugin
* Git
* Network connectivity between VM1 and VM2

VM1 must be able to reach VM2 on ports `8080` and `8081`.

## VM2 — Backend Setup

Clone the repository:

```bash
git clone <repository-url>
cd Materials/cat-dog-app/vm2-app
```

Create the environment file:

```bash
cp .env.example .env
```

Edit it:

```bash
nano .env
```

Example:

```env
POSTGRES_USER=admin
POSTGRES_PASSWORD=change-this-password
POSTGRES_DB=practice
```

### Start the backend

From `vm2-app`:

```bash
sudo docker compose up -d --build
```

Check the containers:

```bash
sudo docker compose ps
```

Check the logs:

```bash
sudo docker compose logs
```

### Verify the APIs

From VM2 or another machine that can reach VM2:

```bash
curl http://192.168.0.20:8080/health
```

Expected:

```text
ok
```

Dog API:

```bash
curl http://192.168.0.20:8081/health
```

Expected:

```text
ok
```

Test the data endpoints:

```bash
curl http://192.168.0.20:8080/cat
```

```bash
curl http://192.168.0.20:8081/dog
```

The APIs return JSON data from PostgreSQL.

## PostgreSQL

PostgreSQL runs inside the Docker Compose network on:

```text
postgres:5432
```

The host port is:

```text
5433
```

Therefore, from the VM host:

```text
localhost:5433
```

can be used to access PostgreSQL.

The database data is persisted outside the container:

```text
/opt/my-db/postgres/data
```

The initialization scripts are mounted from:

```text
postgres/init/
```

### Important PostgreSQL behavior

The SQL initialization scripts are executed when PostgreSQL initializes a new empty data directory.

If PostgreSQL data already exists, changing the SQL initialization files does not automatically recreate the database or tables.

Similarly, changing `POSTGRES_PASSWORD` in `.env` does not change the password of an already initialized PostgreSQL database.

## VM1 — Nginx Setup

The Nginx configuration is located in:

```text
vm1-nginx/
```

The Compose configuration mounts the files into the Nginx container.

The intended deployment location is:

```text
/opt/my-app/nginx
```

The deployment directory contains:

```text
/opt/my-app/nginx/
├── docker-compose.yml
├── nginx.conf
├── conf.d/
├── ssl/
└── www/
```

### TLS Certificate

The repository contains the public certificate:

```text
ssl/cats-dogs.crt
```

The private key must not be committed to Git.

Place the private key on VM1 at:

```text
/opt/my-app/nginx/ssl/cats-dogs.key
```

Make sure the private key is protected with appropriate filesystem permissions.

The Nginx configuration references:

```nginx
ssl_certificate /etc/nginx/ssl/cats-dogs.crt;
ssl_certificate_key /etc/nginx/ssl/cats-dogs.key;
```

## Deploy Nginx

Copy the project to the deployment directory:

```bash
sudo mkdir -p /opt/my-app/nginx
```

Copy the Nginx project files into that directory.

Then:

```bash
cd /opt/my-app/nginx
```

Validate the configuration:

```bash
sudo docker exec nginx nginx -t
```

Expected:

```text
syntax is ok
test is successful
```

Start the container with Compose:

```bash
sudo docker compose up -d
```

Check:

```bash
sudo docker compose ps
```

## DNS / Hosts Configuration

For a local lab, the hostnames can be mapped using `/etc/hosts`.

Example:

```text
192.168.0.10 cats.example.local
192.168.0.10 dogs.example.local
```

Replace `192.168.0.10` with the actual IP address of VM1.

The important point is that both hostnames must resolve to VM1 because VM1 runs Nginx.

## Nginx Reverse Proxy

Nginx selects the backend based on the hostname.

For:

```text
cats.example.local
```

requests are sent to:

```text
http://192.168.0.20:8080
```

For:

```text
dogs.example.local
```

requests are sent to:

```text
http://192.168.0.20:8081
```

For example:

```text
https://cats.example.local/api/cat
        ↓
VM1 Nginx
        ↓
192.168.0.20:8080/cat
```

And:

```text
https://dogs.example.local/api/dog
        ↓
VM1 Nginx
        ↓
192.168.0.20:8081/dog
```

The `/api/` prefix is removed before the request reaches the backend.

## Testing the Complete Application

From VM1:

```bash
curl -k --resolve cats.example.local:443:127.0.0.1 \
  https://cats.example.local/api/cat
```

Dog:

```bash
curl -k --resolve dogs.example.local:443:127.0.0.1 \
  https://dogs.example.local/api/dog
```

You should receive JSON responses from the respective APIs.

## Useful Docker Commands

List running containers:

```bash
sudo docker ps
```

List all containers:

```bash
sudo docker ps -a
```

View Compose services:

```bash
sudo docker compose ps
```

View logs:

```bash
sudo docker compose logs
```

Follow logs:

```bash
sudo docker compose logs -f
```

Stop the application:

```bash
sudo docker compose down
```

Rebuild and start:

```bash
sudo docker compose up -d --build
```

## Troubleshooting

### API is not responding

Check the containers:

```bash
sudo docker compose ps
```

Check API logs:

```bash
sudo docker compose logs cat-api
```

```bash
sudo docker compose logs dog-api
```

Test the APIs directly:

```bash
curl http://192.168.0.20:8080/health
curl http://192.168.0.20:8081/health
```

### PostgreSQL connection errors

Check PostgreSQL:

```bash
sudo docker compose ps postgres
```

Check its logs:

```bash
sudo docker compose logs postgres
```

The APIs should connect to:

```text
postgres:5432
```

inside the Compose network.

They should not use `localhost:5433` for the database connection.

### Nginx configuration errors

Test the configuration:

```bash
sudo docker exec nginx nginx -t
```

View Nginx logs:

```bash
sudo docker logs nginx
```

### Port already in use

Check which process or container is using a port:

```bash
sudo ss -tulpn
```

For Docker containers:

```bash
sudo docker ps
```

## Stopping the Application

On VM1:

```bash
sudo docker compose down
```

On VM2:

```bash
sudo docker compose down
```

PostgreSQL data remains in its persistent storage unless that data directory is deliberately removed.

