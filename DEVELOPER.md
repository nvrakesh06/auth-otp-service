# 🚀 Running the OTP Backend Service

This guide explains how to set up, run, and manage the OTP Backend Service using **Docker & Docker Compose**.

## 📌 Prerequisites
Make sure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/downloads)

## 📂 Step 1: Clone the Repository

Clone your repository on a new system:
```sh
git clone https://github.com/your-username/auth-otp-service.git
cd auth-otp-service
```

## 🐳 Step 2: Ensure Docker is Running
Verify Docker installation:
```sh
docker --version
docker-compose --version
```
📌 If these commands return a version number, Docker is running correctly.

## ⚙️ Step 3: Create the `.env` File
Check if `.env` file exists:
```sh
ls -a  # or `dir /a` on Windows
```
If `.env` is missing, create it:
```sh
cp .env.example .env
```
Then, update `.env` with your database credentials and API keys.

## 🚀 Step 4: Build & Start the Service
Start PostgreSQL, Redis, and the Go backend service:
```sh
docker-compose up --build -d
```

📌 **What this does:**
- Pulls/builds Docker images.
- Starts PostgreSQL & Redis in the background.
- Runs the Go backend service.

## 🛠️ Step 5: Check Running Services
Verify all services are running:
```sh
docker ps
```
✅ **Expected Output:**
```bash
CONTAINER ID   IMAGE           STATUS        PORTS
abcdef123456   postgres:15     Up 10 secs    5432/tcp
123456abcdef   redis:latest    Up 10 secs    6379/tcp
456789abcdef   otp_service     Up 10 secs    8080/tcp
```

If PostgreSQL or Redis isn’t running, restart them:
```sh
docker-compose restart postgres redis
```

## 📜 Step 6: View Logs (If Needed)
Monitor logs:
```sh
docker-compose logs -f go-app
```

✅ **Expected Logs:**
```text
✅ Connected to PostgreSQL successfully!
✅ Connected to Redis successfully!
Fiber v2.0.0 listening on :8080
```

## 🔍 Step 7: Test if Backend is Working
Test backend service:
```sh
curl http://localhost:8080
```

✅ **Expected Response:**
```json
OTP Service is Running!
```

📌 If this fails, check the logs:
```sh
docker-compose logs -f go-app
```

## 🛑 Step 8: Stopping the Service

Stop containers (without deleting data):
```sh
docker-compose down
```

Stop and remove everything (including database volumes):
```sh
docker-compose down --volumes
```

## 🔥 Quick Reference
| Action                      | Command                           |
|-----------------------------|-----------------------------------|
| Start Backend Service       | `docker-compose up --build -d`    |
| Check Running Containers    | `docker ps`                       |
| View Logs                   | `docker-compose logs -f go-app`   |
| Test API Response           | `curl http://localhost:8080`      |
| Stop the Service            | `docker-compose down`             |
| Stop & Remove Everything    | `docker-compose down --volumes`   |