# Profile API Integration Guide

## Overview
The Backend now provides a complete set of endpoints to manage user profiles, including **Basic Info**, **Titles**, **Locations**, and **Interests**.

## 1. API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/v1/users/me/profile` | Get own full profile (includes title, place, interests) | Yes |
| `PATCH` | `/v1/users/me/profile` | Update specific profile fields (bio, title, place, etc.) | Yes |
| `PUT` | `/v1/users/me/full-profile` | **Recommended:** Update Profile AND Interests in one request | Yes |
| `GET` | `/v1/users/:id/profile` | View public profile of another user | No |

## 2. TypeScript Interfaces

Use these interfaces in your Frontend (`types.ts`):

```typescript
// The main User Profile structure
export interface UserProfile {
  id: string;
  username: string;
  first_name: string;
  last_name: string;
  bio?: string;
  avatar_url?: string;
  
  // New Fields
  title?: string;        // e.g. "Senior Engineer"
  place?: string;        // e.g. "San Francisco, USA"
  interests?: Interest[]; // List of interest objects
  
  // Stats
  stats: {
    posts_count: number;
    followers_count: number;
    following_count: number;
  };
}

export interface Interest {
  id: string;
  slug: string;
  name: string;
  category: string;
  icon?: string;
}

// Payload for Composite Update
export interface UpdateFullProfileData {
  first_name?: string;
  last_name?: string;
  username?: string;
  bio?: string;
  
  // New Fields
  title?: string;
  place?: string;
  
  // List of slugs (e.g. ["coding", "design"])
  interest_slugs?: string[];
}
```

## 3. Integration Example (React/Axios)

### Service Layer (`profileService.ts`)

```typescript
import axios from 'axios';
import { UserProfile, UpdateFullProfileData } from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8085/v1';

// Get Current User Profile
export const getMyProfile = async (): Promise<UserProfile> => {
  const response = await axios.get(`${API_URL}/users/me/profile`);
  return response.data;
};

// Update Full Profile (Composite)
export const updateProfile = async (data: UpdateFullProfileData): Promise<UserProfile> => {
  const response = await axios.put(`${API_URL}/users/me/full-profile`, data);
  return response.data;
};
```

### Component Example (`EditProfileForm.tsx`)

```tsx
import { useState } from 'react';
import { updateProfile } from './profileService';

export const EditProfileForm = ({ initialData }) => {
  // State for form fields
  const [formData, setFormData] = useState({
    first_name: initialData.first_name,
    title: initialData.title || '',
    place: initialData.place || '',
    bio: initialData.bio || '',
    interest_slugs: initialData.interests.map(i => i.slug)
  });
  
  const [error, setError] = useState('');

  const handleSave = async () => {
    try {
      const updatedProfile = await updateProfile(formData);
      console.log("Success!", updatedProfile);
      // Update local context/state with new profile
    } catch (err: any) {
      if (err.response?.status === 409) {
        setError("Username is already taken.");
      } else {
        setError("Failed to save profile. Please try again.");
      }
    }
  };

  return (
    <form onSubmit={(e) => { e.preventDefault(); handleSave(); }}>
      {/* Inputs for Title, Place, etc. */}
      {error && <div className="error">{error}</div>}
      <button type="submit">Save Changes</button>
    </form>
  );
};
```

## 4. Key Behaviors

- **Title Handling**: 
  - Send the string (e.g., "Software Engineer").
  - Backend will automatically find the existing Title ID or create a new user-defined Title.
  - Analytics (`used_by_count`) are updated automatically.

- **Place Handling**:
  - Send the string (e.g., "Kota, India").
  - Backend attempts to resolve it to a known City/Country.
  - Usage counts are tracked automatically.

- **Interests**:
  - Send the **Slugs** (e.g., `["coding", "music"]`).
  - The list you send **Replaces** the user's current interests.
