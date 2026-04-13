# Axios Security Review - TravelLink Frontend

## Current Status
- **Axios Version**: 1.15.0 (Latest as of April 2026)
- **Recent Security Advisory**: CVE-2024-39337 (SSRF + URL handling vulnerability)

---

## 🚨 CRITICAL SECURITY ISSUES IDENTIFIED

### 1. **JWT Token Storage in localStorage (XSS Vulnerability)**
**Severity**: HIGH  
**Location**: `frontend/src/services/api.ts` and `frontend/src/stores/auth_store.ts`

**Problem**:
```typescript
localStorage.setItem('access_token', response.token.access_token);
localStorage.setItem('user', JSON.stringify(response.user));
```

✗ **Vulnerabilities**:
- Accessible via `document.cookie` if XSS occurs (JavaScript injection)
- No HttpOnly flag protection
- No Secure flag (if HTTPS used)
- No SameSite protection against CSRF

✓ **Recommendation**: Use secure HTTP-only cookies instead:
```typescript
// Backend should set: Set-Cookie: access_token=...; HttpOnly; Secure; SameSite=Strict
// Frontend reads from cookie automatically (not via JS)
```

---

### 2. **No URL Validation (SSRF Risk)**
**Severity**: HIGH  
**Location**: `frontend/src/services/api.ts` line 9

**Problem**:
```typescript
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
```

✗ **Vulnerabilities**:
- No validation that API_BASE_URL is safe
- Could be manipulated via environment variable injection
- No whitelist of allowed hosts
- Potential SSRF if user input reaches API calls

✓ **Recommendation**:
```typescript
// Validate against whitelist of allowed domains
const ALLOWED_HOSTS = ['localhost:8080', 'api.travellink.com'];
const API_BASE_URL = validateAndGetApiUrl(import.meta.env.VITE_API_URL);

function validateAndGetApiUrl(url: string | undefined): string {
  if (!url) return 'http://localhost:8080/api/v1';
  
  try {
    const urlObj = new URL(url);
    if (!ALLOWED_HOSTS.includes(urlObj.host)) {
      throw new Error(`Unauthorized host: ${urlObj.host}`);
    }
    return url;
  } catch (e) {
    console.error('Invalid API URL configuration', e);
    return 'http://localhost:8080/api/v1';
  }
}
```

---

### 3. **Sensitive Information in Error Messages**
**Severity**: MEDIUM  
**Location**: `frontend/src/stores/auth_store.ts` lines 28-32

**Problem**:
```typescript
error.value = err.response?.data?.message || 'Registration failed';
```

✗ **Vulnerabilities**:
- Backend error messages may expose sensitive data
- Stack traces could leak internal structure
- User enumeration attacks possible

✓ **Recommendation**:
```typescript
const sanitizeErrorMessage = (error: any): string => {
  const status = error.response?.status;
  
  switch (status) {
    case 400: return 'Invalid request. Please check your input.';
    case 401: return 'Authentication failed. Please try again.';
    case 403: return 'Permission denied.';
    case 409: return 'Username or email already in use.';
    case 500: return 'Server error. Please try again later.';
    default: return 'An error occurred. Please try again.';
  }
};
```

---

### 4. **Missing Request Validation & Sanitization**
**Severity**: MEDIUM  
**Location**: Service files (auth_service.ts, plan_service.ts, etc.)

✗ **Vulnerabilities**:
- No input validation before sending to API
- No type validation on responses
- Potential for XSS if response contains unescaped HTML

✓ **Recommendation**:
```typescript
// Validate response types
async function listPlans(params?: ListPlansRequest) {
  const response = await api.get('/plans', { params });
  
  // Validate response structure
  if (!Array.isArray(response.data?.plans)) {
    throw new Error('Invalid server response');
  }
  
  // Validate each plan
  return {
    plans: response.data.plans.map(plan => validatePlan(plan)),
    total: Number(response.data.total),
    page: Number(response.data.page)
  };
}

function validatePlan(plan: any): TravelPlan {
  if (!plan.id || typeof plan.id !== 'string') throw new Error('Invalid plan ID');
  if (!plan.title || typeof plan.title !== 'string') throw new Error('Invalid plan title');
  // ... more validations
  return plan as TravelPlan;
}
```

---

### 5. **No CORS Protection Against Preflight Attacks**
**Severity**: MEDIUM  
**Location**: `frontend/src/services/api.ts` (axios config)

✗ **Vulnerabilities**:
- No explicit CORS handling
- No validation that responses come from expected origin
- Potential for cross-origin data leakage

✓ **Recommendation**:
```typescript
export const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,  // Send cookies with requests
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000,
  responseType: 'json', // Explicitly set expected response type
  validateStatus: (status) => status >= 200 && status < 300, // Only accept 2xx
});

// Validate response origin on critical requests
api.interceptors.response.use(
  (response) => {
    // Validate Content-Type header
    const contentType = response.headers['content-type'];
    if (!contentType?.includes('application/json')) {
      throw new Error('Invalid response content type');
    }
    return response;
  },
  (error) => Promise.reject(error)
);
```

---

### 6. **No Protection Against Man-in-the-Middle (HTTPS)**
**Severity**: HIGH  
**Location**: Environment configuration

✗ **Vulnerabilities**:
- Using `http://` in development but could leak to production
- No certificate pinning
- No HSTS headers checked client-side

✓ **Recommendation**:
```typescript
// Enforce HTTPS in production
const API_BASE_URL = import.meta.env.PROD
  ? 'https://api.travellink.com/api/v1'
  : (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1');

// Add to main.ts
if (import.meta.env.PROD && location.protocol !== 'https:') {
  location.protocol = 'https:';
}
```

---

### 7. **No Rate Limiting at Client Level**
**Severity**: MEDIUM  
**Location**: Service files

✗ **Vulnerabilities**:
- Client could spam API with requests
- No exponential backoff on failures
- Brute force possible (e.g., login attempts)

✓ **Recommendation**: Implement request throttling/debouncing in stores

---

### 8. **Token Refresh Strategy Missing**
**Severity**: MEDIUM  
**Location**: `frontend/src/services/api.ts`

✗ **Vulnerabilities**:
- 1-hour token expiration with no refresh mechanism
- User gets logged out abruptly
- No retry with refresh on 401

✓ **Recommendation**:
```typescript
let isRefreshing = false;
let failedQueue: Array<(token: string) => void> = [];

api.interceptors.response.use(
  response => response,
  async error => {
    const { config, response } = error;
    
    if (response?.status === 401 && !isRefreshing) {
      isRefreshing = true;
      
      try {
        const { data } = await axios.post('/auth/refresh'); // New endpoint needed
        const { access_token } = data;
        localStorage.setItem('access_token', access_token);
        
        api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
        config.headers['Authorization'] = `Bearer ${access_token}`;
        
        failedQueue.forEach(cb => cb(access_token));
        failedQueue = [];
        
        return api(config);
      } finally {
        isRefreshing = false;
      }
    }
    
    return Promise.reject(error);
  }
);
```

---

## 📋 IMMEDIATE ACTIONS (Next Sprint)

### Priority 1 (CRITICAL - Fix this week)
- [ ] Migrate JWT storage to HttpOnly cookies (requires backend change)
- [ ] Add URL validation for API endpoints
- [ ] Sanitize error messages displayed to users
- [ ] Enforce HTTPS in production

### Priority 2 (HIGH - Fix this sprint)
- [ ] Add response validation and type checking
- [ ] Implement request timeout and retry logic
- [ ] Add CORS origin validation
- [ ] Implement token refresh mechanism

### Priority 3 (MEDIUM - Next sprint)
- [ ] Add client-side rate limiting
- [ ] Implement request signing for sensitive operations
- [ ] Add request logging/monitoring
- [ ] Security headers validation

---

## 🔒 Axios Best Practices Applied

✓ **Good Practices Already in Place**:
- Timeout set (30s)
- Authorization header for JWT
- Error handling for 401/403
- Response interceptor for logout on 401

✗ **Missing Best Practices**:
- Request interceptors to add security headers
- Response validation
- Request signing
- Certificate pinning (not applicable in browser)

---

## ⚠️ AXIOS SPECIFIC VULNERABILITIES

### CVE-2024-39337 (SSRF in axios)
- **Affected Versions**: < 1.8.0 (fixed in 1.8.0+)
- **Current Version**: 1.15.0 ✓ (SAFE)
- **Issue**: URL handling in `buildFullPath` could allow SSRF
- **Status**: ✅ PATCHED

---

## 💾 Recommended Changes Summary

```typescript
// BEFORE (Current - UNSAFE)
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  headers: { 'Content-Type': 'application/json' },
  timeout: 30000,
});

localStorage.setItem('access_token', token); // XSS vulnerable

// AFTER (Recommended - SECURE)
const api = axios.create({
  baseURL: validateApiUrl(import.meta.env.VITE_API_URL),
  withCredentials: true, // Use HttpOnly cookies
  headers: { 'Content-Type': 'application/json' },
  timeout: 30000,
  validateStatus: (status) => status >= 200 && status < 300,
});

// Don't store in localStorage - use HttpOnly cookies only
// Backend should set: Set-Cookie: access_token=...; HttpOnly; Secure; SameSite=Strict
```

---

## ✅ Verification Checklist

- [ ] HTTPS enforced in production
- [ ] All tokens stored in HttpOnly cookies only
- [ ] All API URLs whitelisted
- [ ] Error messages sanitized
- [ ] Response types validated
- [ ] CORS headers verified
- [ ] Rate limiting implemented
- [ ] Token refresh working
- [ ] No sensitive data in logs/console
- [ ] Security headers present

