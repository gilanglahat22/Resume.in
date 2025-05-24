#!/bin/bash

echo "=== Checking OAuth Environment Variables ==="
echo ""

echo "GOOGLE_CLIENT_ID: ${GOOGLE_CLIENT_ID:-'NOT SET'}"
echo "GOOGLE_CLIENT_SECRET: ${GOOGLE_CLIENT_SECRET:-'NOT SET'}"  
echo "GOOGLE_REDIRECT_URL: ${GOOGLE_REDIRECT_URL:-'NOT SET'}"
echo ""

echo "=== Expected Configuration ==="
echo "GOOGLE_REDIRECT_URL should be: http://localhost:8080/api/auth/google/callback"
echo ""

if [ -z "$GOOGLE_CLIENT_ID" ]; then
    echo "❌ GOOGLE_CLIENT_ID is not set"
else
    echo "✅ GOOGLE_CLIENT_ID is set"
fi

if [ -z "$GOOGLE_CLIENT_SECRET" ]; then
    echo "❌ GOOGLE_CLIENT_SECRET is not set"  
else
    echo "✅ GOOGLE_CLIENT_SECRET is set"
fi

if [ -z "$GOOGLE_REDIRECT_URL" ]; then
    echo "❌ GOOGLE_REDIRECT_URL is not set"
elif [ "$GOOGLE_REDIRECT_URL" = "http://localhost:8080/api/auth/google/callback" ]; then
    echo "✅ GOOGLE_REDIRECT_URL is correctly set"
else
    echo "⚠️  GOOGLE_REDIRECT_URL is set but might not match: $GOOGLE_REDIRECT_URL"
fi

echo ""
echo "=== Next Steps ==="
echo "1. Set your environment variables in .env file"
echo "2. Configure Google Cloud Console with the exact redirect URI"
echo "3. Restart docker containers: docker-compose down && docker-compose up -d" 