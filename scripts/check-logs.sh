#!/bin/bash

echo "Checking Order Service logs for common errors..."

docker compose logs --tail=200 order-service > logs/order-service.log

echo ""
echo "Detected possible issues:"
grep -iE "error|failed|database|connection|timeout|refused|no such host|missing required" logs/order-service.log || echo "No common error patterns found."

echo ""
echo "Log file saved to logs/order-service.log"