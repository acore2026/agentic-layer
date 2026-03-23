#!/bin/bash

# run_services.sh - Start 6G Agentic Core services and keep them running

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Starting 6G Agentic Core Stack ===${NC}"

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

# 4. Start A-IGW (Translator)
echo -e "${GREEN}[2/3] Starting A-IGW on :18081...${NC}"
go run cmd/igw-fleet/main.go > igw.log 2>&1 &
IGW_PID=$!

# 5. Start AAIHF (Reasoner)
echo -e "${GREEN}[3/3] Starting AAIHF on :18082...${NC}"
go run cmd/aaihf/main.go > aaihf.log 2>&1 &
AAIHF_PID=$!

echo -e "${BLUE}=== All services started ===${NC}"
echo "ACRF PID: $ACRF_PID"
echo "IGW PID: $IGW_PID"
echo "AAIHF PID: $AAIHF_PID"
echo "Logs available in: acrf.log, igw.log, aaihf.log"
echo "Press Ctrl+C to stop all services."

# Wait for Ctrl+C
trap "kill $ACRF_PID $IGW_PID $AAIHF_PID; echo -e '\nServices stopped.'; exit" INT
wait
