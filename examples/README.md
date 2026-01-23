# PowerPro API Examples

This directory contains executable shell scripts that demonstrate how to configure each of the Phase 1 powerlifting programs using the PowerPro API.

## Prerequisites

1. **Running API Server**: Start the PowerPro server:
   ```bash
   go run ./cmd/server
   ```
   The server runs on `http://localhost:8080` by default.

2. **jq** (optional): For pretty-printing JSON responses:
   ```bash
   # macOS
   brew install jq

   # Linux
   apt-get install jq
   ```

## Available Programs

| Script | Program | Demonstrates |
|--------|---------|--------------|
| `starting-strength.sh` | Starting Strength | Linear per-session progression, simple 3x5 |
| `bill-starr-5x5.sh` | Bill Starr 5x5 | Ramping sets, Daily lookups (H/L/M) |
| `wendler-531-bbb.sh` | Wendler 5/3/1 BBB | Weekly lookups, AMRAP sets, cycle progression |
| `sheiko-beginner.sh` | Sheiko Beginner | Multi-week cycles, varied rep schemes |
| `greg-nuckols-high-frequency.sh` | Greg Nuckols High Frequency | Daily max variation, high frequency |

## Usage

Each script is self-contained and will:
1. Create required lifts (if they don't exist)
2. Create prescriptions, days, weeks, and cycles
3. Create progressions appropriate for the program
4. Create the program and link all components
5. Enroll a test user
6. Generate a sample workout

### Running a Single Program

```bash
# Make executable
chmod +x examples/starting-strength.sh

# Run (from project root)
./examples/starting-strength.sh

# Or with pretty JSON output
./examples/starting-strength.sh | jq .
```

### Running All Examples

```bash
# Run all examples
for script in examples/*.sh; do
  echo "=== Running $script ==="
  bash "$script"
  echo ""
done
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `API_BASE_URL` | `http://localhost:8080` | Base URL of the API |
| `ADMIN_USER_ID` | `admin-test-user` | User ID for admin operations |
| `TEST_USER_ID` | `test-user-{program}` | User ID for enrollment |

Example:
```bash
API_BASE_URL=http://localhost:3000 ./examples/starting-strength.sh
```

## Cleanup

To reset the database and run fresh examples:
```bash
# Stop the server, delete the database, restart
rm powerpro.db
go run ./cmd/server
```

## Documentation

For detailed documentation on each program configuration, see:
- `docs/program-configurations/starting-strength.md`
- `docs/program-configurations/bill-starr-5x5.md`
- `docs/program-configurations/wendler-531-bbb.md`
- `docs/program-configurations/sheiko-beginner.md`
- `docs/program-configurations/greg-nuckols-high-frequency.md`

For complete API documentation:
- `docs/openapi.yaml` - OpenAPI 3.0 specification
- `docs/api-reference.md` - Human-readable API reference
- `docs/workflows.md` - Common multi-step workflows
