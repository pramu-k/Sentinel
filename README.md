# Sentinel Monitoring System

Sentinel is a distributed monitoring system that collects system metrics (CPU, Memory, Disk) and service status from Windows servers and visualizes them in a real-time Angular dashboard.

## System Architecture

*   **Sentinel Agent** (Windows Service): Collects metrics and streams them via gRPC.
*   **Sentinel HQ** (Go Server): Receives metrics, stores them in PostgreSQL, and provides a REST API.
*   **Sentinel Dashboard** (Angular): Web interface to view servers, charts, and status.

---

## Prerequisites

*   **Go** (1.21+)
*   **Node.js** (18+) & **npm**
*   **PostgreSQL** (14+)

---

## 1. Database Setup

1.  Ensure PostgreSQL is running.
2.  Create a database named `sentinel`.

```sql
CREATE DATABASE sentinel;
```

*Note: The HQ server will automatically create the necessary tables (`metrics`, `server_status`) on startup.*

---

## 2. Sentinel HQ (Backend)

The HQ server receives data from agents and serves the API.

### Configuration
The HQ is configured via environment variables.
*   `DATABASE_URL`: PostgreSQL connection string. If environment variable is not declared,
default hardcoded connection string in the code will be taken.

### Build & Run
Open a terminal in the project root:

```powershell
# Set connection string (adjust password/host as needed)
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/sentinel?sslmode=disable"

# Run the server
go run cmd/hq/main.go
```

*   **gRPC Server**: Listening on `0.0.0.0:9090`
*   **REST API**: Listening on `0.0.0.0:8080`

---

## 3. Sentinel Agent (Collector)

The Agent runs on the target Windows machine to collect data.

### Configuration
Create a file named `agent-config.json` in the same directory as the agent executable (or project root for dev).

**Example `agent-config.json`**:
```json
{
  "hq_address": "localhost:9090",
  "server_id": "primary-server",
  "collection_interval": "5s",
  "services": [
    "postgres",
    "chrome",
    "sentinel-agent"
  ]
}
```

### Build & Run (Interactive Mode)
```powershell
go run cmd/agent/main.go
```

### Install as Windows Service
To run in the background:
```powershell
# Build the executable first
go build -o sentinel-agent.exe ./cmd/agent

# Install and Start (Requires Admin)
.\sentinel-agent.exe install
.\sentinel-agent.exe start
```

---

## 4. Sentinel Dashboard (Frontend)

The modern web interface.

### Setup
Navigate to the dashboard directory:
```powershell
cd sentinel-dashboard
npm install
```

### Run
```powershell
npm start
```
*   Access the dashboard at **http://localhost:4200**

---

## Troubleshooting

*   **No Data on Dashboard**:
    *   Ensure the Agent is running and connected to HQ (`Connected to HQ` log).
    *   Ensure HQ is connected to the Database.
    *   Check Browser Console (F12) for any API errors.
*   **Agent Connection Failed**:
    *   Verify `hq_address` in `agent-config.json` matches the HQ's gRPC port (default 9090).
