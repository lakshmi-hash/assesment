# User Management App

A small web application with Angular 16 frontend and Go 1.22 backend.

## Setup
1. **Backend:**
   - `cd backend`
   - `go mod tidy`
   - `go run main.go`
   - Swagger: `http://localhost:8080/swagger/index.html`

2. **Frontend:**
   - `cd frontend`
   - `npm install`
   - `npx @angular/cli@16 serve`
   - App: `http://localhost:4200`

## Requirements
- Go 1.22, Node.js, SQLite DLL (in backend folder).