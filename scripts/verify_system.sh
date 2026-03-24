#!/bin/bash

# verify_system.sh - End-to-end verification of the 6G Agentic Core

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Starting 6G Agentic Core Verification ===${NC}"

# 1. Cleanup existing processes
echo -e "Cleaning up existing processes on ports 18080, 18081, 18082, 7233..."
fuser -k 18080/tcp 18081/tcp 18082/tcp 7233/tcp 2>/dev/null || true
pkill -f temporal 2>/dev/null || true

# 2. Check for .env file
if [ ! -f .env ]; then
    echo "Error: .env file not found. Please create one with AGENTIC_GEMINI_API_KEY."
    exit 1
fi

# 3. Start Temporal (New Dependency)
echo -e "${GREEN}[1/4] Starting Temporal Dev Server...${NC}"
/root/.temporalio/bin/temporal server start-dev --ip 127.0.0.1 --log-level info > temporal.log 2>&1 &
TEMPORAL_PID=$!

echo "Waiting for Temporal to be ready..."
TEMPORAL_READY=false
for i in {1..30}; do
    if grep -q "Frontend is now healthy" temporal.log 2>/dev/null; then
        echo -e "${GREEN}Temporal is ready!${NC}"
        TEMPORAL_READY=true
        break
    fi
    sleep 2
done

if [ "$TEMPORAL_READY" = false ]; then
    echo -e "${RED}Error: Temporal failed to start. Check temporal.log for details.${NC}"
    kill $TEMPORAL_PID 2>/dev/null || true
    exit 1
fi

# 4. Start ACRF (Registry)
echo -e "${GREEN}[2/4] Starting ACRF on :18080...${NC}"
go run cmd/acrf/main.go > acrf.log 2>&1 &
ACRF_PID=$!

echo "Waiting for ACRF bootstrap..."
BOOTSTRAP_SUCCESS=false
for i in {1..10}; do
    # Check if acrf.log exists first
    if [ -f acrf.log ] && grep -q "Registered skill: mcp://skill/device/fleet-update" acrf.log; then
        echo -e "${GREEN}ACRF bootstrapped successfully!${NC}"
        BOOTSTRAP_SUCCESS=true
        break
    fi
    sleep 2
done

if [ "$BOOTSTRAP_SUCCESS" = false ]; then
    echo -e "${RED}Error: ACRF bootstrap failed. Check acrf.log for details.${NC}"
    kill $ACRF_PID $TEMPORAL_PID
    exit 1
fi

# 5. Start A-IGW (Translator)
echo -e "${GREEN}[3/4] Starting A-IGW on :18081...${NC}"
export AGENTIC_TEMPORAL_HOST=127.0.0.1:7233
go run cmd/igw-fleet/main.go > igw.log 2>&1 &
IGW_PID=$!
sleep 3

if ! kill -0 $IGW_PID 2>/dev/null; then
    echo -e "${RED}Error: A-IGW failed to start. Check igw.log for details.${NC}"
    kill $ACRF_PID $TEMPORAL_PID
    exit 1
fi

# 6. Start AAIHF (Reasoner)
echo -e "${GREEN}[4/4] Starting AAIHF on :18082...${NC}"
go run cmd/aaihf/main.go > aaihf.log 2>&1 &
AAIHF_PID=$!
sleep 15 

# 7. Send Test Intent
echo -e "${BLUE}=== Sending Test Intent via CURL ===${NC}"
RESPONSE=$(curl -s --noproxy localhost -X POST -H "Content-Type: application/json" \
  -d '{"prompt": "Wake up the fleet for firmware updates", "user_id": "operator-01"}' \
  http://localhost:18082/intent)

echo "$RESPONSE" | python3 -m json.tool

if echo "$RESPONSE" | grep -q "Temporal client is not initialized"; then
    echo -e "${RED}Error: Verification failed because Temporal client was not initialized.${NC}"
    kill $ACRF_PID $IGW_PID $AAIHF_PID $TEMPORAL_PID
    exit 1
fi

echo -e "${BLUE}=== Verification Complete ===${NC}"
echo "Logs available in: acrf.log, igw.log, aaihf.log, temporal.log"

# Cleanup
echo "Stopping services..."
kill $ACRF_PID $IGW_PID $AAIHF_PID $TEMPORAL_PID 2>/dev/null || true
