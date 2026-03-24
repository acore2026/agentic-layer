#!/bin/bash

# verify_system.sh - End-to-end verification of the 6G Agentic Core

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Starting 6G Agentic Core Verification ===${NC}"

# 1. Cleanup existing processes
echo -e "Cleaning up existing processes on ports 18080, 18081, 18082..."
fuser -k 18080/tcp 18081/tcp 18082/tcp 2>/dev/null || true

# 2. Check for .env file
if [ ! -f .env ]; then
    echo "Error: .env file not found. Please create one with AGENTIC_GEMINI_API_KEY."
    exit 1
fi

# 3. Start ACRF (Registry)
echo -e "${GREEN}[1/3] Starting ACRF on :18080...${NC}"
go run cmd/acrf/main.go > acrf.log 2>&1 &
ACRF_PID=$!
sleep 2

RED='\033[0;31m'

# 4. Start A-IGW (Translator)
echo -e "${GREEN}[2/3] Starting A-IGW on :18081...${NC}"
go run cmd/igw-fleet/main.go > igw.log 2>&1 &
IGW_PID=$!
sleep 3

if ! kill -0 $IGW_PID 2>/dev/null; then
    echo -e "${RED}Error: A-IGW failed to start. Check igw.log for details.${NC}"
    kill $ACRF_PID
    exit 1
fi

echo "Waiting for skill registration..."
REGISTRATION_SUCCESS=false
for i in {1..5}; do
    if grep -q "Successfully registered skill" igw.log; then
        echo -e "${GREEN}Skill registration successful!${NC}"
        REGISTRATION_SUCCESS=true
        break
    fi
    sleep 2
done

if [ "$REGISTRATION_SUCCESS" = false ]; then
    echo -e "${RED}Error: Skill registration failed. Check igw.log and acrf.log for details.${NC}"
    kill $ACRF_PID $IGW_PID
    exit 1
fi

# 5. Start AAIHF (Reasoner)
echo -e "${GREEN}[3/3] Starting AAIHF on :18082...${NC}"
# Port and Key are loaded from .env automatically by the binary via godotenv
go run cmd/aaihf/main.go > aaihf.log 2>&1 &
AAIHF_PID=$!
sleep 10 # Give LLM initialization and reasoning some time

# 6. Send Test Intent
echo -e "${BLUE}=== Sending Test Intent via CURL ===${NC}"
RESPONSE=$(curl -s --noproxy localhost -X POST -H "Content-Type: application/json" \
  -d '{"prompt": "Wake up the fleet for firmware updates", "user_id": "operator-01"}' \
  http://localhost:18082/intent)

echo "$RESPONSE" | python3 -m json.tool

if echo "$RESPONSE" | grep -q "Temporal client is not initialized"; then
    echo -e "${RED}Error: Verification failed because Temporal client was not initialized. Please ensure a local Temporal server is running or mock the workflow execution.${NC}"
    kill $ACRF_PID $IGW_PID $AAIHF_PID
    exit 1
fi

echo -e "${BLUE}=== Verification Complete ===${NC}"
echo "Logs available in: acrf.log, igw.log, aaihf.log"

# Cleanup
echo "Stopping services..."
kill $ACRF_PID $IGW_PID $AAIHF_PID
