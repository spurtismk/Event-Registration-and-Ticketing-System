# Setup Guide for Windows

This guide covers everything you need to install and run the Go Event Registration & Ticketing backend on a Windows laptop without Docker or any prior Go setup.

## 1. Install Go

1. Download the Go installer for Windows from: https://go.dev/dl/ (Look for the `.msi` file).
2. Open the downloaded file and follow the standard installation wizard (Next > Next > Install).
3. Open a **new** PowerShell window and verify the installation by running:
   ```powershell
   go version
   ```
   *You should see output like `go version go1.2x.x windows/amd64`.*

## 2. Install PostgreSQL

1. Download the PostgreSQL installer for Windows from EnterpriseDB: https://www.enterprisedb.com/downloads/postgres-postgresql-downloads
2. Run the installer.
3. Keep all the default components checked.
4. **CRITICAL**: When asked for a password, type `postgres` (or remember whatever you set it to).
5. Leave the port as `5432`.
6. Complete the installation.

## 3. Create the Database

Now we need to create the database our app expects (`event_registration`).

1. Search for **pgAdmin 4** in your Windows Start menu and open it (it was installed alongside PostgreSQL).
2. Look at the left sidebar ("Browser"), expand `Servers`, click on `PostgreSQL 16` (or whatever version), and enter the password you set (`postgres`).
3. Right-click on **Databases** > **Create** > **Database...**
4. In the "Database" name field, type: `event_registration`
5. Click **Save**.

## 4. Run the Project

Now that Go and Postgres are ready, you can run the application.

1. Open PowerShell and navigate to the project folder:
   ```powershell
   cd "C:\Users\DELL\OneDrive\Desktop\New folder (5)"
   ```

2. The application needs dependencies installed. Since the folder already has a `go.mod` file, run:
   ```powershell
   go mod tidy
   ```
   *This downloads Gin, GORM, Postgres drivers, JWT, and Crypto packages.*

3. Start the application:
   ```powershell
   go run ./cmd/api
   ```

You should see logs outputting `Starting Server on port 8080...`. The application will automatically create all the necessary database tables (Users, Events, Registrations, etc.) when it starts!

## 5. Test the Application

Leave the server running in that PowerShell window. Open a **second** PowerShell window to test it.

**Register an Admin User:**
```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/auth/register" `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"name": "Admin Test", "email": "admin@test.com", "password": "password123", "role": "ADMIN"}'
```

**Login:**
```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/auth/login" `
  -Headers @{"Content-Type"="application/json"} `
  -Body '{"email": "admin@test.com", "password": "password123"}'
```

The app is now fully functional! Follow the guides in `docs/api_curl_samples.txt` for more advanced commands (note: translating `curl` to `Invoke-RestMethod` in PowerShell, or you can download Postman to test the endpoints visually).
