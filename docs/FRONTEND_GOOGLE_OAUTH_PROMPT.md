# Frontend Integration Prompt: Google OAuth Authentication

## Overview
Integrate Google OAuth authentication with the backend API. Users can sign in with their Google account, which will either create a new account or log into an existing one.

---

## Backend API Base URL
```
http://localhost:8085
```

---

## OAuth Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│ 1. User clicks "Sign in with Google" button                        │
│    → Redirect to: GET /v1/auth/oauth/google                         │
└─────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────┐
│ 2. Backend redirects to Google's OAuth consent screen              │
│    → User signs in with Google account                              │
└─────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────┐
│ 3. Google redirects back to callback URL with auth code            │
│    → GET /v1/auth/oauth/google/callback?code=xxx&state=yyy          │
└─────────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────────┐
│ 4. Backend returns JSON response with tokens and user info         │
│    → Frontend stores tokens and redirects user                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## API Endpoints

### 1. Initiate OAuth Login
```
GET /v1/auth/oauth/google
```

**Behavior**: Redirects the browser to Google's OAuth consent screen.

**Usage**: Open this URL in the browser (not fetch/axios):
```javascript
window.location.href = 'http://localhost:8085/v1/auth/oauth/google';
```

---

### 2. OAuth Callback (handled by backend)
```
GET /v1/auth/oauth/google/callback
```

**Response** (JSON):
```typescript
interface OAuthCallbackResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;      // seconds (900 = 15 minutes)
  is_new_user: boolean;    // true if this is first-time signup
  user: {
    id: string;
    username: string;
    email: string;
  };
}
```

**Example Response**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "is_new_user": true,
  "user": {
    "id": "9a522247-d2b6-4cb2-9848-3f67dcddae1d",
    "username": "",
    "email": "user@gmail.com"
  }
}
```

---

## Frontend Implementation Options

### Option A: Redirect Flow (Recommended for Web)

**Step 1: Create a "Sign in with Google" button**
```tsx
// components/GoogleSignInButton.tsx
export function GoogleSignInButton() {
  const handleGoogleSignIn = () => {
    // Redirect to backend OAuth endpoint
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}/v1/auth/oauth/google`;
  };

  return (
    <button onClick={handleGoogleSignIn} className="google-btn">
      <GoogleIcon />
      Sign in with Google
    </button>
  );
}
```

**Step 2: Create a callback page to handle the response**

Since the backend sends JSON, you need to either:
- Have the backend redirect to frontend with tokens in URL params, OR
- Use a popup flow

**Recommended: Modify backend to redirect to frontend**

If you want the backend to redirect to your frontend after OAuth, you can update the callback to redirect:

```
GET /v1/auth/oauth/google/callback
→ Redirect to: http://localhost:3000/auth/callback?access_token=xxx&refresh_token=xxx&is_new_user=true
```

---

### Option B: Popup Flow

**Step 1: Open OAuth in popup**
```tsx
// hooks/useGoogleAuth.ts
export function useGoogleAuth() {
  const login = () => {
    const popup = window.open(
      `${process.env.NEXT_PUBLIC_API_URL}/v1/auth/oauth/google`,
      'google-oauth',
      'width=500,height=600'
    );

    // Listen for message from popup
    window.addEventListener('message', (event) => {
      if (event.origin !== process.env.NEXT_PUBLIC_API_URL) return;
      
      const { access_token, refresh_token, user, is_new_user } = event.data;
      
      // Store tokens
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', refresh_token);
      
      // Close popup
      popup?.close();
      
      // Handle new user (redirect to complete profile)
      if (is_new_user) {
        window.location.href = '/onboarding';
      } else {
        window.location.href = '/dashboard';
      }
    });
  };

  return { login };
}
```

---

### Option C: Next.js API Route Proxy

**Step 1: Create Next.js API route to initiate**
```typescript
// app/api/auth/google/route.ts
import { redirect } from 'next/navigation';

export async function GET() {
  redirect(`${process.env.API_URL}/v1/auth/oauth/google`);
}
```

**Step 2: Create callback handler**
```typescript
// app/auth/callback/page.tsx
'use client';

import { useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useAuthStore } from '@/stores/authStore';

export default function AuthCallback() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const { login } = useAuthStore();

  useEffect(() => {
    const accessToken = searchParams.get('access_token');
    const refreshToken = searchParams.get('refresh_token');
    const isNewUser = searchParams.get('is_new_user') === 'true';

    if (accessToken && refreshToken) {
      login(accessToken, refreshToken);
      
      if (isNewUser) {
        router.push('/onboarding');
      } else {
        router.push('/dashboard');
      }
    }
  }, [searchParams, login, router]);

  return <div>Signing you in...</div>;
}
```

---

## Storing Tokens

```typescript
// stores/authStore.ts (Zustand)
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import Cookies from 'js-cookie';

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  user: User | null;
  isAuthenticated: boolean;
  login: (accessToken: string, refreshToken: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false,
      
      login: (accessToken, refreshToken) => {
        Cookies.set('auth_token', accessToken, { expires: 1 }); // 1 day
        Cookies.set('refresh_token', refreshToken, { expires: 7 });
        set({ accessToken, refreshToken, isAuthenticated: true });
      },
      
      logout: () => {
        Cookies.remove('auth_token');
        Cookies.remove('refresh_token');
        set({ accessToken: null, refreshToken: null, user: null, isAuthenticated: false });
      },
    }),
    { name: 'auth-storage' }
  )
);
```

---

## Handling New Users

When `is_new_user` is `true`, the user just signed up via Google. They may need to:
1. **Set a username** (required for the app)
2. **Complete their profile** (optional)

Create an onboarding flow:
```typescript
// app/onboarding/page.tsx
export default function OnboardingPage() {
  const [username, setUsername] = useState('');
  
  const handleSubmit = async () => {
    // Update profile with username
    await api.put('/v1/profile', { username });
    router.push('/dashboard');
  };

  return (
    <form onSubmit={handleSubmit}>
      <h1>Complete Your Profile</h1>
      <input
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Choose a username"
      />
      <button type="submit">Continue</button>
    </form>
  );
}
```

---

## Linking Google to Existing Account

If a user is already logged in and wants to link their Google account:

```typescript
// POST /v1/auth/oauth/google/link
// Requires: Authorization header with access token
// Body: { "code": "authorization_code_from_google" }

const linkGoogle = async (code: string) => {
  const response = await api.post('/v1/auth/oauth/google/link', { code });
  // Response: { message: "Account linked successfully", provider: "google", provider_email: "..." }
};
```

---

## Unlinking Google Account

```typescript
// DELETE /v1/auth/oauth/google
// Requires: Authorization header

const unlinkGoogle = async () => {
  await api.delete('/v1/auth/oauth/google');
};
```

---

## Error Handling

```typescript
// Common OAuth errors
const errors = {
  'Invalid oauth state': 'Session expired. Please try again.',
  'provider google not configured': 'Google sign-in is not available.',
  'failed to exchange code': 'Authentication failed. Please try again.',
};
```

---

## UI Component Example

```tsx
// components/SocialLoginButtons.tsx
export function SocialLoginButtons() {
  return (
    <div className="social-buttons">
      <button
        onClick={() => window.location.href = '/api/auth/google'}
        className="btn-google"
      >
        <svg>...</svg>
        Continue with Google
      </button>
      
      {/* Add more providers as needed */}
    </div>
  );
}
```

---

## Security Considerations

1. **State Parameter**: Backend handles CSRF protection via state cookie
2. **HTTPS**: Use HTTPS in production
3. **Token Storage**: Store in httpOnly cookies when possible
4. **Refresh Tokens**: Use the refresh endpoint to get new access tokens
