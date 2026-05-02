# Cheat CLI - Quick Start Guide

The **Cheat CLI** is a command-line tool for executing course-related tasks. It's lightweight, single-user focused, and perfect for automation and scripting.

---

## Table of Contents

- [What is Cheat CLI?](#what-is-cheat-cli)
- [Building the CLI](#building-the-cli)
- [Running the CLI](#running-the-cli)
- [Authentication](#authentication)
- [Usage Examples](#usage-examples)
- [Troubleshooting](#troubleshooting)

---

## What is Cheat CLI?

The Cheat CLI is a single-user command-line interface that allows you to:

- ✅ Authenticate with email and password
- ✅ Execute course-related operations
- ✅ View course information
- ✅ Process course enrollments
- ✅ Generate progress reports
- ✅ Run in automated scripts and cron jobs

**CLI vs Service:** The CLI is perfect for one-time executions or scheduled tasks. For multi-user concurrent operations, use the **Cheat Service** instead.

---

## Building the CLI

### Option 1: Build for Your Current Platform

```bash
cd /path/to/cheat-master
go build -o cheat-cli ./cmd/cli/main.go
```

**Output:** `cheat-cli` (macOS/Linux) or `cheat-cli.exe` (Windows)

### Option 2: Build for Linux (Ubuntu 24.04, x64)

```bash
cd /path/to/cheat-master
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cheat-cli ./cmd/cli/main.go
```

**Output:** `cheat-cli` (static Linux binary)

### Option 3: Use the Build Script (Recommended)

The automated build script handles all compilation flags:

```bash
./build-azure.sh
# When prompted, select: 1 (CLI Mode Only)
```

**Output Location:** `dist/cheat-cli`

#### Build Script Features:

- ✅ Automatic version detection (Go 1.20+)
- ✅ Pre-flight validation checks
- ✅ Binary architecture verification
- ✅ SHA256 checksum generation
- ✅ Cross-platform support (macOS → Linux)

---

## Running the CLI

### Basic Execution

```bash
./cheat-cli
```

### What to Expect

1. **Email Prompt:**

   ```
   Email:
   ```

   Enter your registered email address

2. **Password Prompt:**

   ```
   Password:
   ```

   Enter your password (input is hidden)

3. **Processing:**
   The CLI authenticates and executes course operations

---

## Authentication

### Credentials

- **Email:** Your registered email address
- **Password:** Your account password

### First-Time Users

If this is your first time using the CLI, you may need to:

1. Register your credentials with the service:

   ```bash
   # Use the Service mode to register
   ./cheat-service
   # Select option 1: Register new user
   ```

2. Then use the CLI with those credentials:
   ```bash
   ./cheat-cli
   ```

### Credential Storage

Credentials are:

- ✅ Stored securely in the service's credentials.json
- ✅ Used only for authentication
- ✅ Never sent over unencrypted networks (when deployed properly)

---

## Usage Examples

### Example 1: Local Execution

```bash
$ ./cheat-cli
Email: user@example.com
Password:
# Processing course operations...
# Results displayed on console
```

### Example 2: Scripted Execution (piping credentials)

```bash
# Create a credentials file (EXAMPLE ONLY - NOT RECOMMENDED FOR PRODUCTION)
echo "user@example.com" > /tmp/creds.txt
echo "password123" >> /tmp/creds.txt

# Pipe to CLI
cat /tmp/creds.txt | ./cheat-cli
```

### Example 3: Scheduled Execution (cron)

```bash
# Every day at 9 AM
0 9 * * * /path/to/cheat-cli < /secure/credentials.txt >> /var/log/cheat-cli.log 2>&1
```

### Example 4: In a Shell Script

```bash
#!/bin/bash

CLI_BINARY="./cheat-cli"
EMAIL="user@example.com"
PASSWORD="password123"

# Execute with input
(echo "$EMAIL"; echo "$PASSWORD") | $CLI_BINARY
```

---

## Troubleshooting

### ❌ "command not found: cheat-cli"

**Problem:** Binary not in PATH or doesn't exist

**Solution:**

```bash
# Build it first
./build-azure.sh
# Select option 1 (CLI Mode)

# Or run with explicit path
./dist/cheat-cli
```

### ❌ "permission denied"

**Problem:** Binary doesn't have execute permissions

**Solution:**

```bash
chmod +x cheat-cli
./cheat-cli
```

### ❌ "authentication failed"

**Problem:** Email or password is incorrect

**Solution:**

1. Verify your credentials
2. Register with the Service mode if not registered:
   ```bash
   ./cheat-service
   # Select: 1 (Register new user)
   ```
3. Try again with correct credentials

### ❌ "connection refused" / "no such file"

**Problem:** Service credentials or data files not accessible

**Solution:**

1. Ensure you're in the correct directory:
   ```bash
   pwd
   # Should show: /path/to/cheat-master
   ```
2. Verify service was initialized:
   ```bash
   ls -la internal/models/
   ls -la configs/
   ```

### ❌ Binary crashes immediately

**Problem:** Architecture mismatch (wrong binary for platform)

**Solution:**

```bash
# Check your system architecture
uname -m
# x86_64 = amd64
# arm64 = arm64

# Rebuild with correct target
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cheat-cli ./cmd/cli/main.go
```

---

## Performance Tips

### Optimize for Batch Operations

If running multiple CLI invocations:

```bash
# ❌ SLOW: Multiple separate invocations
./cheat-cli
./cheat-cli
./cheat-cli

# ✅ FAST: Use the Service mode instead
# Service keeps connections open and processes jobs concurrently
./cheat-service
```

### Reduce Output

```bash
# Suppress output to speed up execution
./cheat-cli > /dev/null 2>&1
```

---

## Deployment

### To Azure Linux VM

```bash
# 1. Build on local machine
./build-azure.sh
# Select: 1 (CLI Mode)

# 2. Transfer to Azure VM
scp dist/cheat-cli azureuser@your-vm-ip:~/

# 3. Connect and run
ssh azureuser@your-vm-ip
chmod +x cheat-cli
./cheat-cli
```

### In Docker

```dockerfile
FROM ubuntu:24.04

# Copy binary
COPY dist/cheat-cli /app/cheat-cli
RUN chmod +x /app/cheat-cli

ENTRYPOINT ["/app/cheat-cli"]
```

### In Kubernetes

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cheat-cli-job
spec:
  containers:
    - name: cli
      image: cheat:latest
      imagePullPolicy: IfNotPresent
  restartPolicy: Never
```

---

## Comparison: CLI vs Service

| Feature         | CLI                 | Service                  |
| --------------- | ------------------- | ------------------------ |
| **Users**       | Single              | Multiple                 |
| **Persistence** | Temp                | Permanent (JSON)         |
| **Concurrency** | Sequential          | Concurrent (4 workers)   |
| **Best For**    | Automation, scripts | Long-running, multi-user |
| **Memory**      | Low                 | Medium                   |
| **Setup**       | Simple              | More complex             |
| **Cost**        | Low                 | Medium                   |

---

## Next Steps

- 📖 Read [BUILD_GUIDE.md](BUILD_GUIDE.md) for advanced build options
- 🚀 Deploy to Azure VM using [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md)
- 📊 Monitor jobs with [Service Mode](./cmd/service/README.md)
- 🔒 Review security best practices in [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md)

---

## Support

**Having issues?** Check:

1. ✅ [Troubleshooting](#troubleshooting) section above
2. ✅ [BUILD_GUIDE.md](BUILD_GUIDE.md) - Detailed build troubleshooting
3. ✅ Binary architecture: `file cheat-cli`
4. ✅ System compatibility: `uname -m` and `uname -s`

---

**Last Updated:** May 2, 2026  
**Version:** Compatible with Go 1.20+  
**License:** [Your Project License]
