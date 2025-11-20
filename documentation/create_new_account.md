# Signup & Onboarding API Design

This document defines the API design for a multi-step onboarding flow supporting:

* Name & DOB
* Username & Avatar
* Email / Mobile Signup
* Password / Social Login
* OTP / Email Verification
* Interests Setup
* Token issuance & refresh

---

## **1. POST /v1/onboard/basic**

Submit first name, last name, and date of birth.

### Request Body

```json
{
  "temp_id": "optional-uuid",
  "first_name": "John",
  "last_name": "Doe",
  "dob": "1999-05-12"
}
```

### Response

**201 Created**

```json
{
  "temp_id": "generated-uuid",
  "message": "basic_saved"
}
```

---

## **2. POST /v1/onboard/profile**

Submit username + optional avatar upload.

### Request Body

```json
{
  "temp_id": "uuid",
  "username": "john123",
  "avatar_upload": {
    "filename": "avatar.jpg",
    "content_type": "image/jpeg"
  }
}
```

### Response

```json
{
  "temp_id": "uuid",
  "username_available": true,
  "avatar_url": "https://cdn.example.com/avatars/uuid.jpg"
}
```

---

## **3. POST /v1/auth/signup**

Start signup using email, mobile, or social login.

### Request Body (email/password)

```json
{
  "temp_id": "uuid",
  "method": "email",
  "email": "user@example.com",
  "password": "Secret123!"
}
```

### Request Body (mobile)

```json
{
  "temp_id": "uuid",
  "method": "mobile",
  "phone": "9876543210",
  "country_code": "+91",
  "password": "Secret123!"
}
```

### Request Body (social)

```json
{
  "temp_id": "uuid",
  "method": "google",
  "social_token": "google-id-token"
}
```

### Response

```json
{
  "user_id": "uuid",
  "verification_needed": true,
  "verification_channel": "email"
}
```

---

## **4. POST /v1/auth/verify**

Verify OTP or email token.

### Request Body

```json
{
  "user_id": "uuid",
  "method": "email",
  "code": "123456"
}
```

### Response

```json
{
  "access_token": "<jwt>",
  "refresh_token": "<refresh_jwt>",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "username": "john123",
    "email": "user@example.com",
    "phone": "+919876543210"
  }
}
```

---

## **5. POST /v1/onboard/interests** (Auth Required)

Save user interests.

### Request Body

```json
{
  "interests": ["JEE Mains", "NEET", "SSC"]
}
```

### Response

```json
{
  "message": "interests_saved"
}
```

---

## **6. POST /v1/auth/login**

Login using email or phone.

### Request Body

```json
{
  "identifier": "user@example.com",
  "password": "Secret123!"
}
```

### Response

```json
{
  "access_token": "<jwt>",
  "refresh_token": "<refresh_jwt>",
  "expires_in": 900
}
```

---

## **7. POST /v1/auth/refresh**

Refresh the access token.

### Request Body

```json
{
  "refresh_token": "<refresh_jwt>"
}
```

### Response

```json
{
  "access_token": "<new_jwt>",
  "refresh_token": "<new_refresh_jwt>",
  "expires_in": 900
}
```

---

## **8. POST /v1/auth/logout** (Auth Required)

Revoke refresh token.

### Request Body

```json
{
  "refresh_token": "<refresh_jwt>"
}
```

### Response

```json
{
  "message": "logged_out"
}
```

---

## **9. POST /v1/auth/resend**

Resend OTP/email verification.

### Request Body

```json
{
  "user_id": "uuid",
  "method": "email"
}
```

### Response

```json
{
  "message": "code_sent"
}
```

---

## **10. POST /v1/upload/avatar** (Auth Required)

Supports avatar upload.

### Multipart Request

```
file: image/jpeg, image/png
```

### Response

```json
{
  "avatar_url": "https://cdn.example.com/avatars/uuid.jpg"
}
```

---

## **Data Models (Postgres)**

### Users

```sql
users (
  id uuid primary key,
  username text unique,
  email text unique,
  email_verified bool,
  phone text unique,
  phone_verified bool,
  first_name text,
  last_name text,
  dob date,
  password_hash text,
  avatar_url text,
  created_at timestamptz,
  updated_at timestamptz
)
```

### Refresh Tokens

```sql
refresh_tokens (
  id uuid primary key,
  user_id uuid,
  token_hash text,
  issued_at timestamptz,
  expires_at timestamptz,
  revoked bool
)
```

### Interests

```sql
user_interests (
  user_id uuid,
  interest_slug text
)
```

---

## **Error Response Format**

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Invalid email format",
    "fields": { "email": "invalid_format" }
  }
}
```

---

## **Common Error Codes**

* `INVALID_INPUT`
* `EMAIL_ALREADY_EXISTS`
* `PHONE_ALREADY_EXISTS`
* `USERNAME_TAKEN`
* `WEAK_PASSWORD`
* `OTP_INVALID`
* `OTP_EXPIRED`
* `TOO_MANY_OTP_ATTEMPTS`
* `REFRESH_TOKEN_INVALID`
* `REFRESH_TOKEN_REUSED`

---

This `.md` file contains the full signup + onboarding API design. Let me know if you want an **OpenAPI (Swagger) version**, a **diagram**, or **database migrations**.
