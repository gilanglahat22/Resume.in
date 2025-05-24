#!/bin/bash

echo "🔍 Testing OAuth Setup for redirect_uri_mismatch fix"
echo "=================================================="
echo ""

echo "1️⃣ Testing registration endpoint..."
RESPONSE=$(curl -s http://localhost:8080/api/auth/google/register)
echo "Response: $RESPONSE"
echo ""

if [[ $RESPONSE == *"auth_url"* ]]; then
    echo "✅ Registration endpoint is working"
    
    # Extract auth_url
    AUTH_URL=$(echo $RESPONSE | grep -o '"auth_url":"[^"]*' | cut -d'"' -f4)
    echo "🔗 Auth URL: $AUTH_URL"
    echo ""
    
    # Extract redirect_uri from the URL
    REDIRECT_URI=$(echo $AUTH_URL | grep -o 'redirect_uri=[^&]*' | cut -d'=' -f2)
    DECODED_URI=$(python3 -c "import urllib.parse; print(urllib.parse.unquote('$REDIRECT_URI'))" 2>/dev/null || echo "$REDIRECT_URI")
    echo "📍 Redirect URI being sent: $DECODED_URI"
    echo ""
    
    echo "2️⃣ What you need to configure in Google Cloud Console:"
    echo "   Go to: https://console.cloud.google.com"
    echo "   Navigate to: APIs & Services > Credentials"
    echo "   Edit your OAuth 2.0 Client ID"
    echo "   Add this EXACT redirect URI:"
    echo "   ➡️  $DECODED_URI"
    echo ""
    
    echo "3️⃣ Required scopes in OAuth consent screen:"
    echo "   - ../auth/userinfo.email"
    echo "   - ../auth/userinfo.profile"
    echo "   - ../auth/user.birthday.read"
    echo "   - ../auth/user.gender.read"
    echo "   - ../auth/user.phonenumbers.read"
    echo ""
    
    echo "4️⃣ Test the OAuth URL:"
    echo "   Copy this URL and paste in browser:"
    echo "   $AUTH_URL"
    echo ""
    echo "   If you see 'redirect_uri_mismatch', the redirect URI is not configured correctly."
    
else
    echo "❌ Registration endpoint failed"
    echo "   Check if backend is running: docker-compose up -d backend postgres"
    echo "   Check environment variables are set"
fi

echo ""
echo "📋 Quick checklist:"
echo "   □ Google Cloud Console OAuth client configured"
echo "   □ Exact redirect URI added: http://localhost:8080/api/auth/google/callback"
echo "   □ OAuth consent screen with required scopes"
echo "   □ Test users added (for External user type)"
echo "   □ Environment variables set in .env file"
echo "   □ Backend restarted after environment changes" 