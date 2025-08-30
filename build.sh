#!/bin/bash
set -e

PURPLE='\033[0;35m'
NC='\033[0m'

if ! command -v docker >/dev/null 2>&1; then
        echo "Docker is not installed. Please install it first."
        exit 1
fi

if ! command -v bun >/dev/null 2>&1; then
        echo "Bun is not installed. Please install it first."
        exit 1
fi

echo -e "${PURPLE}[VID.SH] === Starting backend build ===${NC}"
cd backend
docker build  --quiet -t vidsh-backend . &

echo -e "${PURPLE}[VID.SH] === Starting frontend build ===${NC}"
cd ../frontend
bun i >/dev/null
bun run build >/dev/null

wait

echo -e "${PURPLE}[VID.SH] === Build finished! ===${NC}"
echo -e "${PURPLE}[VID.SH] === If you don't know what do to next check out the README.md file ===${NC}"