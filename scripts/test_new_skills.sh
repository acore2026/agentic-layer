#!/bin/bash

# test_new_skills.sh - Test the new network skills end-to-end (Build Mode)

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Testing New Network Skills (Build Mode) ===${NC}"

# 1. Cleanup
fuser -k 18080/tcp 18081/tcp 18082/tcp 7233/tcp 2>/dev/null || true
pkill -f temporal || true
rm -f acrf.log igw.log aaihf.log temporal.log

export AGENTIC_TEMPORAL_HOST=127.0.0.1:7233

# 2. Start Temporal
echo "Starting Temporal..."
/root/.temporalio/bin/temporal server start-dev --ip 127.0.0.1 > temporal.log 2>&1 &
TEMPORAL_PID=$!

echo "Waiting for Temporal to be ready..."
for i in {1..20}; do
    if grep -q "Frontend is now healthy" temporal.log; then
        echo "Temporal is ready!"
        break
    fi
    sleep 2
done

# 3. Build
echo "Building binaries..."
go build -o acrf-bin cmd/acrf/main.go
go build -o igw-bin cmd/igw-fleet/main.go
go build -o aaihf-bin cmd/aaihf/main.go

# 4. Start
echo "Starting ACRF..."
./acrf-bin > acrf.log 2>&1 &
ACRF_PID=$!

echo "Starting A-IGW..."
./igw-bin > igw.log 2>&1 &
IGW_PID=$!

echo "Starting AAIHF..."
export AGENTIC_USE_MOCK_AGENT=true
./aaihf-bin > aaihf.log 2>&1 &
AAIHF_PID=$!

echo "Waiting for services to start..."
sleep 15

# 5. Test intents
echo -e "${GREEN}Testing Turbo Mode (QoS)...${NC}"
curl -s --noproxy localhost -X POST -H "Content-Type: application/json" \
  -d '{"prompt": "Enable Turbo Mode for Gaming Session.", "user_id": "operator-01"}' \
  http://localhost:18082/intent | python3 -m json.tool
sleep 2

echo -e "\n${GREEN}Testing Resiliency (Reliability)...${NC}"
curl -s --noproxy localhost -X POST -H "Content-Type: application/json" \
  -d '{"prompt": "Ensure Zero-Interruption for V2X Feed", "user_id": "operator-01"}' \
  http://localhost:18082/intent | python3 -m json.tool
sleep 2

echo -e "\n${GREEN}Testing Secure Drone Corridor (Edge)...${NC}"
curl -s --noproxy localhost -X POST -H "Content-Type: application/json" \
  -d '{"prompt": "Secure Drone Corridor.", "user_id": "operator-01"}' \
  http://localhost:18082/intent | python3 -m json.tool
sleep 5

echo -e "\n${BLUE}=== Checking IGW Logs for Signaling Sequences ===${NC}"
grep -E "Starting Translation Sequence|Mock 3GPP" igw.log

# Cleanup
echo -e "\n${BLUE}=== Testing Complete ===${NC}"
kill $ACRF_PID $IGW_PID $AAIHF_PID $TEMPORAL_PID 2>/dev/null || true
rm -f acrf-bin igw-bin aaihf-bin
