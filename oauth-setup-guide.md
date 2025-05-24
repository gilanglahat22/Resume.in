# Google OAuth Setup Guide - REDIRECT_URI_MISMATCH FIX

## 🚨 **CRITICAL**: Current redirect_uri_mismatch Error Solution

Your application is sending this exact redirect URI to Google:
```
http://localhost:8080/api/auth/google/callback
```

**You MUST configure this exact URI in Google Cloud Console.**

## Required Google Cloud Console Setup

### 1. Create OAuth 2.0 Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project or create a new one
3. Navigate to **APIs & Services** > **Credentials**
4. Click **+ CREATE CREDENTIALS** > **OAuth 2.0 Client IDs**

### 2. Configure OAuth Consent Screen

1. Go to **APIs & Services** > **OAuth consent screen**
2. Choose **External** user type (for testing)
3. Fill in required fields:
   - **App name**: Resume.in
   - **User support email**: your-email@domain.com
   - **Developer contact email**: your-email@domain.com
4. Add these **Scopes** (REQUIRED for registration):
   - `../auth/userinfo.email`
   - `../auth/userinfo.profile` 
   - `../auth/user.birthday.read`
   - `../auth/user.gender.read`
   - `../auth/user.phonenumbers.read`
5. Add **Test users** (your email addresses for testing)

### 3. 🎯 Configure OAuth Client - **THIS IS THE KEY STEP**

1. **Application type**: Web application
2. **Name**: Resume.in OAuth Client
3. **Authorized redirect URIs**: Click **"+ ADD URI"** and add **EXACTLY**:
   ```
   http://localhost:8080/api/auth/google/callback
   ```
   
   **⚠️ CRITICAL REQUIREMENTS:**
   - NO trailing slash
   - NO extra spaces
   - Exact case matching
   - Must be exactly: `http://localhost:8080/api/auth/google/callback`

4. **Click "SAVE"**

### 4. Get Credentials

After creation, copy:
- **Client ID**: `xxxxx.apps.googleusercontent.com`
- **Client Secret**: `xxxxxx`

## Environment Configuration

Create a `.env` file in the root directory with YOUR actual credentials:

```env
# Authentication Configuration - REPLACE WITH YOUR ACTUAL VALUES
GOOGLE_CLIENT_ID=323216826330-dmc9q6gj8pu5e3mfcbnag7radt912kjb.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-actual-client-secret-here
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# Other Configuration
SERVER_PORT=8080
ENVIRONMENT=development
ALLOW_ORIGINS="http://localhost:3000,http://localhost:4200"
JWT_SECRET=your-super-secret-jwt-key-change-in-production
FRONTEND_URL=http://localhost:4200

# Database Configuration
POSTGRES_USER=resumeuser
POSTGRES_PASSWORD=resumepassword
POSTGRES_DB=resumedb
```

## OAuth Flow Explanation

### Unified Callback Approach
- **Single Redirect URI**: Both login and registration use the same callback
- **State-Based Detection**: Registration uses `register_` prefix in state parameter
- **Enhanced Scopes**: Registration requests additional permissions (birthday, gender, phone)

### How It Works
1. **Login**: `/api/auth/google/login` → Google OAuth → `/api/auth/google/callback`
2. **Registration**: `/api/auth/google/register` → Google OAuth → `/api/auth/google/callback`
3. Backend detects flow type from state parameter prefix

## Testing the Setup

1. **Set your environment variables in `.env` file**

2. **Restart the backend**:
   ```bash
   docker-compose down
   docker-compose up -d backend postgres
   ```

3. **Test registration endpoint**:
   ```bash
   curl http://localhost:8080/api/auth/google/register
   ```

4. **Should return**:
   ```json
   {
     "auth_url": "https://accounts.google.com/o/oauth2/auth?..."
   }
   ```

5. **Test the auth_url in browser** - should NOT show "redirect_uri_mismatch" error

## Troubleshooting redirect_uri_mismatch

### Error: "redirect_uri_mismatch"
This error occurs when the redirect URI sent to Google doesn't match any authorized redirect URIs in your Google Cloud Console.

**Step-by-step fix:**

1. **Check what your app is sending**:
   ```bash
   curl http://localhost:8080/api/auth/google/register
   ```
   
2. **Look for `redirect_uri=` in the auth_url** (URL-decode it)

3. **Go to Google Cloud Console** → APIs & Services → Credentials

4. **Edit your OAuth 2.0 Client ID**

5. **In "Authorized redirect URIs"**, make sure you have EXACTLY**:
   ```
   http://localhost:8080/api/auth/google/callback
   ```

6. **Save and wait 5-10 minutes** for changes to propagate

7. **Test again**

### Other Common Issues

### "Invalid scope" Error  
- **Cause**: Missing scopes in OAuth consent screen
- **Solution**: Add all required scopes in Google Console

### "Unauthorized_client" Error
- **Cause**: Missing or incorrect client credentials
- **Solution**: Verify `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET`

### "Access denied" Error
- **Cause**: User not added as test user (for External app type)
- **Solution**: Add your email as test user in OAuth consent screen

## Production Setup

For production, update:
- **Authorized redirect URIs**: `https://yourdomain.com/api/auth/google/callback`
- **GOOGLE_REDIRECT_URL**: `https://yourdomain.com/api/auth/google/callback`
- **OAuth consent screen**: Submit for verification if needed
- **User type**: Change to Internal if using Google Workspace

## Quick Checklist

- [ ] Google Cloud Console project created
- [ ] OAuth consent screen configured with all required scopes
- [ ] OAuth 2.0 Client ID created with exact redirect URI: `http://localhost:8080/api/auth/google/callback`
- [ ] `.env` file created with your actual GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET
- [ ] Backend restarted after setting environment variables
- [ ] Test users added (if using External user type)
- [ ] Your email added as test user

## Security Notes

- Keep client secret secure and never expose in frontend code
- Use HTTPS in production
- Implement proper CSRF protection with state parameters
- Regularly rotate JWT secrets
- Monitor OAuth usage in Google Console 