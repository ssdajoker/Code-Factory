# Mode Specifications

**Version:** 1.0.0  
**Status:** Draft  
**Last Updated:** 2026-01-07  
**Owner:** Code-Factory Core Team

---

## Overview

The Spec-Driven Software Factory operates in four distinct modes, each designed for a specific phase of the software development lifecycle. This document provides detailed specifications for each mode.

---

## 1. INTAKE Mode

### 1.1 Purpose

Capture project vision and requirements, then transform them into structured, actionable specifications using LLM assistance.

### 1.2 User Journey

```
Start INTAKE → Describe Vision → LLM Generates Spec → Review & Edit → Save → Commit (Optional)
```

### 1.3 Interface Design

**Main Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         INTAKE MODE                                  ║
║                  Capture Your Project Vision                         ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 📝 Describe your feature or project:                                │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ I want to build a user authentication system with:                  │
│ - Email/password login                                              │
│ - OAuth (Google, GitHub)                                            │
│ - JWT tokens                                                         │
│ - Password reset flow                                               │
│ - Rate limiting                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│ [Ctrl+Enter to generate spec | Ctrl+C to cancel]                    │
└──────────────────────────────────────────────────────────────────────┘

💡 Tip: Be as detailed as possible. Include requirements, constraints,
        and any specific technologies you want to use.
```

**Generating Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         INTAKE MODE                                  ║
║                  Generating Specification...                         ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 🤖 AI is analyzing your requirements...                             │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ ⏳ Step 1/4: Analyzing requirements... ✓                            │
│ ⏳ Step 2/4: Structuring specification... ⏳                         │
│ ⏸️  Step 3/4: Adding technical details...                           │
│ ⏸️  Step 4/4: Generating examples...                                │
│                                                                      │
│ Estimated time: 15 seconds                                           │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

**Review & Edit Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         INTAKE MODE                                  ║
║                  Review Generated Specification                      ║
╚══════════════════════════════════════════════════════════════════════╝

┌─────────────────────────────────┬────────────────────────────────────┐
│ 📄 Specification Preview        │ ✏️  Edit Mode                      │
├─────────────────────────────────┼────────────────────────────────────┤
│                                 │                                    │
│ # User Authentication System    │ # User Authentication System       │
│                                 │                                    │
│ ## Overview                     │ ## Overview                        │
│ A secure authentication system  │ A secure authentication system     │
│ supporting multiple login       │ supporting multiple login          │
│ methods...                      │ methods...                         │
│                                 │                                    │
│ ## Requirements                 │ ## Requirements                    │
│ 1. Email/password login         │ 1. Email/password login            │
│ 2. OAuth integration            │ 2. OAuth integration               │
│    - Google                     │    - Google                        │
│    - GitHub                     │    - GitHub                        │
│ 3. JWT token management         │ 3. JWT token management            │
│ ...                             │ ...                                │
│                                 │                                    │
│ [Tab to switch | ↑↓ to scroll] │ [Ctrl+S to save | Esc to cancel]   │
└─────────────────────────────────┴────────────────────────────────────┘

Actions: [S]ave | [E]dit | [R]egenerate | [C]ancel
```

### 1.4 Workflow Steps

#### Step 1: Input Capture
- **Input Method:** Multi-line text editor
- **Validation:** Minimum 50 characters
- **Enhancements:**
  - Syntax highlighting for markdown
  - Auto-save to temp file every 30 seconds
  - Restore from temp on crash

#### Step 2: LLM Processing
- **Prompt Template:** `prompts/intake.tmpl`
- **Context:** Project name, existing specs, tech stack
- **Output Format:** Structured markdown
- **Streaming:** Display generation in real-time

#### Step 3: Review & Edit
- **Split View:** Preview (left) + Edit (right)
- **Live Preview:** Update preview as user edits
- **Validation:** Check for required sections
- **Suggestions:** AI-powered improvements

#### Step 4: Save
- **Filename:** Auto-generated from title (e.g., `user_authentication.md`)
- **Location:** `/contracts/` or user-specified subdirectory
- **Metadata:** Add frontmatter (date, author, version)
- **Git:** Optionally commit with message

### 1.5 LLM Prompt Template

```markdown
You are a technical specification writer. Generate a comprehensive, 
actionable specification document based on the following requirements.

## User Requirements:
{user_input}

## Project Context:
- Project Name: {project_name}
- Tech Stack: {tech_stack}
- Existing Specs: {existing_specs}

## Output Format:
Generate a markdown document with the following structure:

# [Feature Name]

## Overview
Brief description of the feature (2-3 sentences)

## Requirements
### Functional Requirements
1. Requirement 1
2. Requirement 2
...

### Non-Functional Requirements
1. Performance: ...
2. Security: ...
3. Scalability: ...

## Technical Design
### Architecture
Describe the high-level architecture

### Components
List and describe key components

### Data Models
Define data structures and schemas

### API Endpoints (if applicable)
List API endpoints with methods and parameters

## Implementation Plan
1. Phase 1: ...
2. Phase 2: ...
...

## Testing Strategy
- Unit tests
- Integration tests
- E2E tests

## Security Considerations
List security concerns and mitigations

## Performance Considerations
Expected performance characteristics

## Dependencies
External libraries, services, or APIs

## Open Questions
List any unresolved questions or decisions

---

Generate the specification now:
```

### 1.6 Output Example

```markdown
# User Authentication System

**Version:** 1.0.0
**Status:** Draft
**Created:** 2026-01-07
**Author:** ssdajoker

## Overview

A secure, scalable authentication system supporting multiple login methods 
including email/password and OAuth providers (Google, GitHub). The system 
uses JWT tokens for session management and includes comprehensive security 
features like rate limiting and password reset flows.

## Requirements

### Functional Requirements

1. **Email/Password Authentication**
   - User registration with email verification
   - Secure password hashing (bcrypt, min 10 rounds)
   - Login with email and password
   - Remember me functionality (optional)

2. **OAuth Integration**
   - Google OAuth 2.0
   - GitHub OAuth 2.0
   - Automatic account linking for existing emails

3. **JWT Token Management**
   - Access tokens (15-minute expiry)
   - Refresh tokens (7-day expiry)
   - Token rotation on refresh
   - Blacklist for revoked tokens

4. **Password Reset Flow**
   - Request reset via email
   - Secure reset token (1-hour expiry)
   - Password strength validation
   - Email notification on successful reset

5. **Rate Limiting**
   - Login attempts: 5 per 15 minutes per IP
   - Password reset: 3 per hour per email
   - Token refresh: 10 per minute per user

### Non-Functional Requirements

1. **Performance**
   - Login response time: < 200ms (p95)
   - Token validation: < 10ms (p99)
   - Support 1000 concurrent users

2. **Security**
   - OWASP Top 10 compliance
   - HTTPS only
   - Secure cookie flags (HttpOnly, Secure, SameSite)
   - CSRF protection

3. **Scalability**
   - Horizontal scaling support
   - Stateless authentication (JWT)
   - Redis for rate limiting and token blacklist

## Technical Design

### Architecture

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   Client    │─────▶│  API Server │─────▶│  Database   │
│  (Browser)  │◀─────│   (Go)      │◀─────│ (PostgreSQL)│
└─────────────┘      └─────────────┘      └─────────────┘
                            │
                            ▼
                     ┌─────────────┐
                     │    Redis    │
                     │ (Rate Limit)│
                     └─────────────┘
```

### Components

1. **Auth Service** (`internal/auth`)
   - User registration and login
   - Password hashing and verification
   - Token generation and validation

2. **OAuth Service** (`internal/oauth`)
   - OAuth flow handling
   - Provider-specific implementations
   - Account linking logic

3. **Token Service** (`internal/token`)
   - JWT generation and parsing
   - Token refresh logic
   - Blacklist management

4. **Rate Limiter** (`internal/ratelimit`)
   - IP-based rate limiting
   - User-based rate limiting
   - Redis-backed storage

### Data Models

**User:**
```go
type User struct {
    ID            uuid.UUID
    Email         string
    PasswordHash  string
    EmailVerified bool
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

**OAuthAccount:**
```go
type OAuthAccount struct {
    ID         uuid.UUID
    UserID     uuid.UUID
    Provider   string // "google", "github"
    ProviderID string
    CreatedAt  time.Time
}
```

**RefreshToken:**
```go
type RefreshToken struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Token     string
    ExpiresAt time.Time
    CreatedAt time.Time
}
```

### API Endpoints

**POST /auth/register**
- Body: `{ "email": "user@example.com", "password": "..." }`
- Response: `{ "user_id": "...", "message": "Verification email sent" }`

**POST /auth/login**
- Body: `{ "email": "user@example.com", "password": "..." }`
- Response: `{ "access_token": "...", "refresh_token": "..." }`

**POST /auth/refresh**
- Body: `{ "refresh_token": "..." }`
- Response: `{ "access_token": "...", "refresh_token": "..." }`

**POST /auth/logout**
- Headers: `Authorization: Bearer <access_token>`
- Response: `{ "message": "Logged out successfully" }`

**GET /auth/oauth/{provider}**
- Redirects to OAuth provider

**GET /auth/oauth/{provider}/callback**
- Query: `code=...&state=...`
- Response: `{ "access_token": "...", "refresh_token": "..." }`

**POST /auth/password-reset/request**
- Body: `{ "email": "user@example.com" }`
- Response: `{ "message": "Reset email sent" }`

**POST /auth/password-reset/confirm**
- Body: `{ "token": "...", "new_password": "..." }`
- Response: `{ "message": "Password reset successfully" }`

## Implementation Plan

1. **Phase 1: Core Authentication (Week 1)**
   - User model and database schema
   - Email/password registration
   - Email/password login
   - JWT token generation

2. **Phase 2: OAuth Integration (Week 2)**
   - Google OAuth implementation
   - GitHub OAuth implementation
   - Account linking logic

3. **Phase 3: Security Features (Week 3)**
   - Rate limiting
   - Password reset flow
   - Token refresh and rotation
   - Token blacklist

4. **Phase 4: Testing & Hardening (Week 4)**
   - Unit tests (>80% coverage)
   - Integration tests
   - Security audit
   - Performance testing

## Testing Strategy

**Unit Tests:**
- Password hashing and verification
- JWT generation and parsing
- Rate limiter logic
- OAuth flow handlers

**Integration Tests:**
- Full registration flow
- Full login flow
- OAuth flow (mocked providers)
- Password reset flow

**E2E Tests:**
- User registration and login via UI
- OAuth login via UI
- Password reset via UI

**Security Tests:**
- SQL injection attempts
- XSS attempts
- CSRF attacks
- Brute force login attempts

## Security Considerations

1. **Password Security**
   - Bcrypt with cost factor 12
   - Minimum 8 characters, require uppercase, lowercase, number
   - Check against common password lists

2. **Token Security**
   - Short-lived access tokens (15 minutes)
   - Secure refresh token storage
   - Token rotation on refresh
   - Blacklist for revoked tokens

3. **OAuth Security**
   - Validate state parameter
   - Verify redirect URI
   - Use PKCE for mobile apps

4. **Rate Limiting**
   - Prevent brute force attacks
   - Prevent account enumeration
   - Prevent DoS attacks

## Performance Considerations

- **Database Indexes:** Email (unique), UserID (foreign keys)
- **Caching:** User sessions in Redis (optional)
- **Connection Pooling:** PostgreSQL connection pool (max 100)
- **Horizontal Scaling:** Stateless design allows multiple instances

## Dependencies

- **Go Libraries:**
  - `github.com/golang-jwt/jwt/v5` - JWT handling
  - `golang.org/x/crypto/bcrypt` - Password hashing
  - `golang.org/x/oauth2` - OAuth 2.0 client
  - `github.com/go-redis/redis/v9` - Redis client
  - `github.com/lib/pq` - PostgreSQL driver

- **External Services:**
  - PostgreSQL 14+
  - Redis 6+
  - SMTP server (for emails)

## Open Questions

1. Should we support 2FA (TOTP)?
2. Should we implement social login (Facebook, Twitter)?
3. What's the password reset token expiry time? (Proposed: 1 hour)
4. Should we log all authentication events for audit?

---

**End of Specification**
```

### 1.7 Error Handling

**LLM Unavailable:**
- Offer template-based spec creation
- Allow manual writing with structure hints

**Invalid Input:**
- Show validation errors inline
- Suggest improvements

**Save Failure:**
- Retry with exponential backoff
- Offer alternative save location
- Keep content in memory

---

## 2. REVIEW Mode

### 2.1 Purpose

Compare existing code against specifications to identify compliance issues, deviations, and areas for improvement.

### 2.2 User Journey

```
Start REVIEW → Select Spec → Scan Code → LLM Analysis → View Results → Generate Report → (Optional) Create Issues
```

### 2.3 Interface Design

**Spec Selection Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         REVIEW MODE                                  ║
║                  Check Code Against Specifications                   ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 📋 Select specification to review:                                  │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ ↓ user_authentication.md (Modified 2 days ago)                      │
│   api_endpoints.md (Modified 1 week ago)                            │
│   database_schema.md (Modified 2 weeks ago)                         │
│   deployment_process.md (Modified 1 month ago)                      │
│                                                                      │
│ [↑↓ to navigate | Enter to select | A to review all | / to search] │
└──────────────────────────────────────────────────────────────────────┘

Options:
[A]ll Specs | [M]odified Only | [S]earch | [C]ancel
```

**Scanning Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         REVIEW MODE                                  ║
║                  Scanning Codebase...                                ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 🔍 Analyzing: user_authentication.md                                │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ Step 1/4: Scanning codebase... ✓                                    │
│   Found 23 relevant files                                           │
│                                                                      │
│ Step 2/4: Loading specifications... ✓                               │
│   Loaded 1 specification                                            │
│                                                                      │
│ Step 3/4: Analyzing code vs spec... ⏳                              │
│   Progress: [████████████░░░░░░░░] 65% (15/23 files)               │
│   Current: internal/auth/service.go                                 │
│                                                                      │
│ Step 4/4: Generating report... ⏸️                                   │
│                                                                      │
│ Estimated time remaining: 45 seconds                                 │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

**Results Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         REVIEW MODE                                  ║
║                  Review Results                                      ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 📊 Compliance Summary                                               │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ Specification: user_authentication.md                               │
│ Files Analyzed: 23                                                   │
│ Overall Compliance: 78% ⚠️                                           │
│                                                                      │
│ ✅ Compliant: 18 items                                              │
│ ⚠️  Warnings: 4 items                                               │
│ ❌ Issues: 3 items                                                  │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ 🔴 Critical Issues (3)                                              │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ 1. ❌ Missing rate limiting on /auth/login endpoint                │
│    File: internal/auth/handlers.go:45                               │
│    Spec: Section 2.1.5 - Rate Limiting                              │
│    Impact: High - Security vulnerability                            │
│    [View Details] [Suggest Fix] [Create Issue]                      │
│                                                                      │
│ 2. ❌ Password hashing uses bcrypt cost 10 (spec requires 12)      │
│    File: internal/auth/password.go:23                               │
│    Spec: Section 4.1 - Password Security                            │
│    Impact: Medium - Reduced security                                │
│    [View Details] [Suggest Fix] [Create Issue]                      │
│                                                                      │
│ 3. ❌ Missing CSRF protection on OAuth callback                    │
│    File: internal/oauth/handlers.go:78                              │
│    Spec: Section 4.3 - OAuth Security                               │
│    Impact: High - Security vulnerability                            │
│    [View Details] [Suggest Fix] [Create Issue]                      │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ [↑↓ to navigate | Enter for details | R for report | Q to quit]    │
└──────────────────────────────────────────────────────────────────────┘
```

**Issue Detail Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         REVIEW MODE                                  ║
║                  Issue Details                                       ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ ❌ Missing rate limiting on /auth/login endpoint                    │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ File: internal/auth/handlers.go:45                                   │
│ Severity: High                                                       │
│ Category: Security                                                   │
│                                                                      │
│ Specification Requirement:                                           │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │ Rate Limiting:                                                 │  │
│ │ - Login attempts: 5 per 15 minutes per IP                      │  │
│ │ - Password reset: 3 per hour per email                         │  │
│ │ - Token refresh: 10 per minute per user                        │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ Current Implementation:                                              │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │ func (h *Handler) Login(w http.ResponseWriter, r *http.Request) { │
│ │     // No rate limiting implemented                            │  │
│ │     var req LoginRequest                                       │  │
│ │     json.NewDecoder(r.Body).Decode(&req)                       │  │
│ │     ...                                                        │  │
│ │ }                                                              │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ Suggested Fix:                                                       │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │ func (h *Handler) Login(w http.ResponseWriter, r *http.Request) { │
│ │     // Add rate limiting                                       │  │
│ │     if !h.rateLimiter.Allow(r.RemoteAddr, "login", 5, 15*time.Minute) { │
│ │         http.Error(w, "Too many requests", http.StatusTooManyRequests) │
│ │         return                                                 │  │
│ │     }                                                          │  │
│ │     var req LoginRequest                                       │  │
│ │     json.NewDecoder(r.Body).Decode(&req)                       │  │
│ │     ...                                                        │  │
│ │ }                                                              │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ Actions:                                                             │
│ [C]opy Fix | [O]pen in Editor | [I]ssue on GitHub | [N]ext | [B]ack│
└──────────────────────────────────────────────────────────────────────┘
```

### 2.4 Workflow Steps

#### Step 1: Spec Selection
- List all specs in `/contracts/`
- Show last modified date
- Allow filtering by date, name, or status
- Option to review all specs at once

#### Step 2: Code Scanning
- Use git to find relevant files
- Filter by file extensions (configurable)
- Exclude vendor/, node_modules/, etc.
- Parallel file reading for performance

#### Step 3: LLM Analysis
- Send spec + code to LLM
- Ask for compliance check
- Request specific issues and suggestions
- Stream results for immediate feedback

#### Step 4: Results Display
- Categorize issues: Critical, Warning, Info
- Show compliance percentage
- Highlight specific code locations
- Provide actionable suggestions

#### Step 5: Report Generation
- Generate markdown report
- Save to `/reports/review_YYYY-MM-DD.md`
- Include summary, issues, and suggestions
- Optionally commit to git

#### Step 6: Issue Creation (Optional)
- Create GitHub issues for critical problems
- Link to spec and code location
- Include suggested fix
- Assign to team members

### 2.5 LLM Prompt Template

```markdown
You are a code reviewer. Compare the following code against the specification
and identify compliance issues, deviations, and areas for improvement.

## Specification:
{spec_content}

## Code Files:
{code_files}

## Analysis Instructions:
1. Check if all requirements from the spec are implemented
2. Identify any deviations from the spec
3. Look for security issues mentioned in the spec
4. Check for performance considerations from the spec
5. Verify data models match the spec
6. Verify API endpoints match the spec

## Output Format:
For each issue found, provide:

### Issue: [Brief description]
- **Severity:** Critical | High | Medium | Low
- **Category:** Security | Performance | Functionality | Style
- **File:** [file path]:[line number]
- **Spec Section:** [section reference]
- **Current Implementation:** [code snippet]
- **Expected Implementation:** [what the spec requires]
- **Suggested Fix:** [code snippet or description]
- **Impact:** [description of impact]

Also provide:
- **Overall Compliance:** [percentage]
- **Summary:** [brief summary of findings]
- **Compliant Items:** [list of things that are correct]

---

Begin analysis:
```

### 2.6 Report Example

```markdown
# Code Review Report

**Date:** 2026-01-07  
**Specification:** user_authentication.md  
**Reviewer:** Factory AI  
**Overall Compliance:** 78%

---

## Executive Summary

The user authentication system implementation is mostly compliant with the 
specification, but has 3 critical security issues that need immediate attention.
The core functionality is implemented correctly, but some security features 
specified in the document are missing.

**Key Findings:**
- ✅ 18 requirements fully implemented
- ⚠️ 4 requirements partially implemented
- ❌ 3 requirements not implemented (security-critical)

---

## Compliance Breakdown

| Category | Compliant | Warnings | Issues |
|----------|-----------|----------|--------|
| Functional Requirements | 5/5 | 0 | 0 |
| Security Requirements | 2/5 | 1 | 2 |
| Performance Requirements | 3/3 | 0 | 0 |
| API Endpoints | 7/8 | 1 | 0 |
| Data Models | 3/3 | 0 | 0 |

---

## Critical Issues (3)

### 1. ❌ Missing Rate Limiting on Login Endpoint

**Severity:** High  
**Category:** Security  
**File:** `internal/auth/handlers.go:45`  
**Spec Section:** 2.1.5 - Rate Limiting

**Specification Requirement:**
```
Rate Limiting:
- Login attempts: 5 per 15 minutes per IP
- Password reset: 3 per hour per email
- Token refresh: 10 per minute per user
```

**Current Implementation:**
```go
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    // No rate limiting implemented
    var req LoginRequest
    json.NewDecoder(r.Body).Decode(&req)
    ...
}
```

**Suggested Fix:**
```go
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    // Add rate limiting
    if !h.rateLimiter.Allow(r.RemoteAddr, "login", 5, 15*time.Minute) {
        http.Error(w, "Too many requests", http.StatusTooManyRequests)
        return
    }
    var req LoginRequest
    json.NewDecoder(r.Body).Decode(&req)
    ...
}
```

**Impact:** Without rate limiting, the login endpoint is vulnerable to brute 
force attacks. An attacker could attempt unlimited password guesses.

---

### 2. ❌ Password Hashing Uses Bcrypt Cost 10 (Spec Requires 12)

**Severity:** Medium  
**Category:** Security  
**File:** `internal/auth/password.go:23`  
**Spec Section:** 4.1 - Password Security

**Specification Requirement:**
```
Password Security:
- Bcrypt with cost factor 12
- Minimum 8 characters, require uppercase, lowercase, number
```

**Current Implementation:**
```go
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    return string(hash), err
}
```

**Suggested Fix:**
```go
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}
```

**Impact:** Lower bcrypt cost makes password hashes easier to crack. Cost 12 
is the current recommended minimum for security.

---

### 3. ❌ Missing CSRF Protection on OAuth Callback

**Severity:** High  
**Category:** Security  
**File:** `internal/oauth/handlers.go:78`  
**Spec Section:** 4.3 - OAuth Security

**Specification Requirement:**
```
OAuth Security:
- Validate state parameter
- Verify redirect URI
- Use PKCE for mobile apps
```

**Current Implementation:**
```go
func (h *Handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    // Missing state validation
    token, err := h.oauth.Exchange(r.Context(), code)
    ...
}
```

**Suggested Fix:**
```go
func (h *Handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")
    
    // Validate state parameter
    if !h.oauth.ValidateState(state) {
        http.Error(w, "Invalid state", http.StatusBadRequest)
        return
    }
    
    token, err := h.oauth.Exchange(r.Context(), code)
    ...
}
```

**Impact:** Without state validation, the OAuth flow is vulnerable to CSRF 
attacks. An attacker could trick a user into authorizing their account.

---

## Warnings (4)

### 1. ⚠️ Token Expiry Time Differs from Spec

**Severity:** Low  
**Category:** Functionality  
**File:** `internal/token/jwt.go:34`  
**Spec Section:** 2.1.3 - JWT Token Management

**Issue:** Access tokens expire in 30 minutes, but spec specifies 15 minutes.

**Suggested Action:** Update token expiry to match spec, or update spec if 
30 minutes is intentional.

---

### 2. ⚠️ Missing Email Notification on Password Reset

**Severity:** Low  
**Category:** Functionality  
**File:** `internal/auth/password_reset.go:56`  
**Spec Section:** 2.1.4 - Password Reset Flow

**Issue:** Password reset succeeds but doesn't send confirmation email.

**Suggested Action:** Add email notification after successful password reset.

---

### 3. ⚠️ Database Connection Pool Size Not Configured

**Severity:** Low  
**Category:** Performance  
**File:** `internal/database/db.go:12`  
**Spec Section:** 5 - Performance Considerations

**Issue:** Connection pool uses default size, spec recommends max 100.

**Suggested Action:** Configure connection pool: `db.SetMaxOpenConns(100)`

---

### 4. ⚠️ Missing Index on Users.Email

**Severity:** Medium  
**Category:** Performance  
**File:** `migrations/001_create_users.sql:5`  
**Spec Section:** 5 - Performance Considerations

**Issue:** No index on email column, which is used for lookups.

**Suggested Action:** Add index: `CREATE UNIQUE INDEX idx_users_email ON users(email);`

---

## Compliant Items (18)

✅ User registration with email verification  
✅ Secure password hashing (bcrypt)  
✅ Login with email and password  
✅ Google OAuth integration  
✅ GitHub OAuth integration  
✅ JWT token generation  
✅ Refresh token rotation  
✅ Token blacklist for revoked tokens  
✅ Password strength validation  
✅ User data model matches spec  
✅ OAuth account data model matches spec  
✅ Refresh token data model matches spec  
✅ POST /auth/register endpoint  
✅ POST /auth/login endpoint  
✅ POST /auth/refresh endpoint  
✅ POST /auth/logout endpoint  
✅ GET /auth/oauth/{provider} endpoint  
✅ POST /auth/password-reset/request endpoint  

---

## Recommendations

1. **Immediate Actions (Critical Issues):**
   - Implement rate limiting on all authentication endpoints
   - Increase bcrypt cost to 12
   - Add CSRF protection to OAuth callback

2. **Short-Term Actions (Warnings):**
   - Align token expiry times with spec
   - Add email notification on password reset
   - Configure database connection pool
   - Add database indexes

3. **Long-Term Improvements:**
   - Consider adding 2FA support
   - Implement audit logging for all auth events
   - Add monitoring and alerting for failed login attempts

---

## Next Steps

1. Create GitHub issues for critical problems
2. Assign issues to team members
3. Schedule security review after fixes
4. Re-run Factory review after implementation

---

**Generated by Factory v1.0.0**
```

---

## 3. CHANGE_ORDER Mode

### 3.1 Purpose

Track specification drift over time by detecting code changes that deviate from 
specifications, and manage the process of updating specs or reverting changes.

### 3.2 User Journey

```
Start CHANGE_ORDER → Detect Changes → Analyze Drift → Review → Approve/Reject → Update Spec or Create Issue
```

### 3.3 Interface Design

**Change Detection Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                      CHANGE_ORDER MODE                               ║
║                  Track Specification Drift                           ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 🔍 Detecting changes since last review...                           │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ Comparing: HEAD vs last-review (7 days ago)                         │
│                                                                      │
│ Files changed: 12                                                    │
│ Lines added: +234                                                    │
│ Lines removed: -89                                                   │
│                                                                      │
│ Analyzing changes against specifications...                         │
│ Progress: [████████████████░░░░] 80% (10/12 files)                 │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

**Drift Analysis Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                      CHANGE_ORDER MODE                               ║
║                  Specification Drift Analysis                        ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 📊 Drift Summary                                                    │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ Total Changes: 12 files                                              │
│ Spec-Compliant: 8 files ✅                                          │
│ Spec Drift: 4 files ⚠️                                              │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ 🔴 Changes Requiring Attention (4)                                  │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ 1. ⚠️ New endpoint added: POST /auth/2fa/enable                    │
│    File: internal/auth/handlers.go (+45 lines)                      │
│    Spec: user_authentication.md (no mention of 2FA)                 │
│    Type: Feature Addition                                           │
│    [View Details] [Approve] [Reject]                                │
│                                                                      │
│ 2. ⚠️ Token expiry changed: 15min → 30min                          │
│    File: internal/token/jwt.go (1 line)                             │
│    Spec: user_authentication.md (specifies 15min)                   │
│    Type: Behavior Change                                            │
│    [View Details] [Approve] [Reject]                                │
│                                                                      │
│ 3. ⚠️ Database schema modified: added 'phone' column               │
│    File: migrations/003_add_phone.sql (+5 lines)                    │
│    Spec: user_authentication.md (no phone field in User model)      │
│    Type: Data Model Change                                          │
│    [View Details] [Approve] [Reject]                                │
│                                                                      │
│ 4. ⚠️ Removed rate limiting from /auth/refresh                     │
│    File: internal/auth/handlers.go (-8 lines)                       │
│    Spec: user_authentication.md (requires rate limiting)            │
│    Type: Spec Violation                                             │
│    [View Details] [Approve] [Reject]                                │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ [↑↓ to navigate | Enter for details | A to approve all | Q to quit]│
└──────────────────────────────────────────────────────────────────────┘
```

**Change Detail Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                      CHANGE_ORDER MODE                               ║
║                  Change Order Details                                ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ ⚠️ New endpoint added: POST /auth/2fa/enable                        │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ File: internal/auth/handlers.go                                      │
│ Commit: a3f5b2c - "Add 2FA support" (2 days ago)                     │
│ Author: ssdajoker                                                    │
│ Spec: user_authentication.md                                         │
│                                                                      │
│ Change Type: Feature Addition                                        │
│ Drift Severity: Medium                                               │
│                                                                      │
│ Code Changes:                                                        │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │ + func (h *Handler) Enable2FA(w http.ResponseWriter, r *http.Request) { │
│ │ +     userID := r.Context().Value("user_id").(uuid.UUID)       │  │
│ │ +     secret, err := totp.Generate(totp.GenerateOpts{          │  │
│ │ +         Issuer:      "MyApp",                                │  │
│ │ +         AccountName: userID.String(),                        │  │
│ │ +     })                                                       │  │
│ │ +     if err != nil {                                          │  │
│ │ +         http.Error(w, err.Error(), http.StatusInternalServerError) │
│ │ +         return                                               │  │
│ │ +     }                                                        │  │
│ │ +     // Store secret in database                             │  │
│ │ +     h.db.Save2FASecret(userID, secret.Secret())             │  │
│ │ +     json.NewEncoder(w).Encode(map[string]string{            │  │
│ │ +         "secret": secret.Secret(),                           │  │
│ │ +         "qr_code": secret.URL(),                             │  │
│ │ +     })                                                       │  │
│ │ + }                                                            │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ Specification Status:                                                │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │ ❌ Not mentioned in specification                              │  │
│ │                                                                │  │
│ │ The specification has an "Open Questions" section that asks:  │  │
│ │ "Should we support 2FA (TOTP)?"                                │  │
│ │                                                                │  │
│ │ This change appears to answer that question with "yes".        │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ AI Analysis:                                                         │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │ This is a feature addition that was anticipated in the spec    │  │
│ │ (mentioned in Open Questions). The implementation looks        │  │
│ │ reasonable, but the specification should be updated to:        │  │
│ │                                                                │  │
│ │ 1. Document the 2FA feature in Requirements section            │  │
│ │ 2. Add API endpoint documentation                              │  │
│ │ 3. Update data model to include 2FA secret field               │  │
│ │ 4. Add security considerations for 2FA                         │  │
│ │                                                                │  │
│ │ Recommendation: APPROVE and update specification               │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ Actions:                                                             │
│ [A]pprove & Update Spec | [R]eject & Create Issue | [S]kip | [B]ack│
└──────────────────────────────────────────────────────────────────────┘
```

**Approval Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                      CHANGE_ORDER MODE                               ║
║                  Approve Change & Update Specification               ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ ✅ Approving change: POST /auth/2fa/enable                          │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ The specification will be updated with the following changes:        │
│                                                                      │
│ 📝 user_authentication.md                                           │
│                                                                      │
│ Section: 2.1 Functional Requirements                                 │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │ + 6. **Two-Factor Authentication (2FA)**                       │  │
│ │ +    - TOTP-based 2FA                                          │  │
│ │ +    - QR code generation for easy setup                       │  │
│ │ +    - Backup codes for account recovery                       │  │
│ │ +    - Optional enforcement per user                           │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ Section: 3.3 Data Models                                             │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │   type User struct {                                           │  │
│ │       ID            uuid.UUID                                  │  │
│ │       Email         string                                     │  │
│ │       PasswordHash  string                                     │  │
│ │       EmailVerified bool                                       │  │
│ │ +     TwoFactorSecret string                                   │  │
│ │ +     TwoFactorEnabled bool                                    │  │
│ │       CreatedAt     time.Time                                  │  │
│ │       UpdatedAt     time.Time                                  │  │
│ │   }                                                            │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ Section: 3.4 API Endpoints                                           │
│ ┌────────────────────────────────────────────────────────────────┐  │
│ │ + **POST /auth/2fa/enable**                                    │  │
│ │ + - Headers: `Authorization: Bearer <access_token>`            │  │
│ │ + - Response: `{ "secret": "...", "qr_code": "..." }`          │  │
│ │ +                                                              │  │
│ │ + **POST /auth/2fa/verify**                                    │  │
│ │ + - Body: `{ "code": "123456" }`                               │  │
│ │ + - Response: `{ "success": true }`                            │  │
│ │ +                                                              │  │
│ │ + **POST /auth/2fa/disable**                                   │  │
│ │ + - Headers: `Authorization: Bearer <access_token>`            │  │
│ │ + - Body: `{ "code": "123456" }`                               │  │
│ │ + - Response: `{ "success": true }`                            │  │
│ └────────────────────────────────────────────────────────────────┘  │
│                                                                      │
│ Change Order Record:                                                 │
│ - ID: CO-001                                                         │
│ - Date: 2026-01-07                                                   │
│ - Type: Feature Addition                                             │
│ - Status: Approved                                                   │
│ - Commit: a3f5b2c                                                    │
│                                                                      │
│ [C]onfirm | [E]dit Changes | [C]ancel                               │
└──────────────────────────────────────────────────────────────────────┘
```

### 3.4 Workflow Steps

#### Step 1: Change Detection
- Use git diff to find changes since last review
- Filter by relevant file types
- Group changes by commit or file
- Calculate change statistics

#### Step 2: Spec Mapping
- For each changed file, find related specs
- Use file path patterns and content analysis
- LLM can help identify relevant specs

#### Step 3: Drift Analysis
- Compare changes against spec requirements
- Classify drift type: Addition, Modification, Deletion, Violation
- Assess severity: Low, Medium, High, Critical
- Generate explanation and recommendation

#### Step 4: Review & Decision
- Present changes to user
- Show code diff, spec reference, and AI analysis
- User decides: Approve, Reject, or Skip
- Record decision in change order log

#### Step 5: Action Execution
- **If Approved:** Update specification with changes
- **If Rejected:** Create GitHub issue to revert or fix
- **If Skipped:** Log for future review

#### Step 6: Change Order Documentation
- Create change order document in `/reports/change_orders/`
- Include: ID, date, type, status, commits, decisions
- Track history of spec evolution

### 3.5 LLM Prompt Template

```markdown
You are a specification analyst. Analyze the following code changes and 
determine if they comply with the specification or represent drift.

## Specification:
{spec_content}

## Code Changes (git diff):
{git_diff}

## Analysis Instructions:
1. Identify what changed in the code
2. Find the relevant section in the specification
3. Determine if the change is:
   - Compliant: Implements something from the spec
   - Addition: Adds new functionality not in spec
   - Modification: Changes existing behavior from spec
   - Violation: Contradicts or removes something from spec

4. Assess severity:
   - Low: Minor change, no impact
   - Medium: Notable change, should update spec
   - High: Significant change, requires review
   - Critical: Breaking change or security issue

5. Provide recommendation: Approve, Reject, or Needs Discussion

## Output Format:

### Change Summary
[Brief description of what changed]

### Drift Analysis
- **Type:** Compliant | Addition | Modification | Violation
- **Severity:** Low | Medium | High | Critical
- **Spec Section:** [reference to spec section]
- **Current Spec:** [what the spec says]
- **New Behavior:** [what the code now does]

### Impact Assessment
[Description of impact on system, users, security, performance]

### Recommendation
**Action:** Approve | Reject | Needs Discussion

**Reasoning:** [explanation of recommendation]

**If Approve:** [suggested spec updates]
**If Reject:** [suggested fix or revert]

---

Begin analysis:
```

### 3.6 Change Order Document Example

```markdown
# Change Order CO-001

**Date:** 2026-01-07  
**Type:** Feature Addition  
**Status:** Approved  
**Specification:** user_authentication.md  
**Commits:** a3f5b2c, b4e6d3f

---

## Summary

Added two-factor authentication (2FA) support using TOTP. This feature was 
anticipated in the specification's "Open Questions" section but not fully 
specified.

---

## Changes

### Code Changes

**Files Modified:**
- `internal/auth/handlers.go` (+45 lines)
- `internal/auth/2fa.go` (+120 lines, new file)
- `migrations/003_add_2fa.sql` (+10 lines, new file)

**New Endpoints:**
- `POST /auth/2fa/enable` - Enable 2FA for user
- `POST /auth/2fa/verify` - Verify 2FA code
- `POST /auth/2fa/disable` - Disable 2FA for user

**Data Model Changes:**
- Added `two_factor_secret` column to `users` table
- Added `two_factor_enabled` column to `users` table

### Specification Changes

**Sections Updated:**
- 2.1 Functional Requirements - Added 2FA requirement
- 3.3 Data Models - Updated User model
- 3.4 API Endpoints - Added 2FA endpoints
- 4.1 Security Considerations - Added 2FA security notes
- 6 Implementation Plan - Added 2FA to Phase 2

---

## Drift Analysis

**Type:** Feature Addition  
**Severity:** Medium  
**Spec Status:** Anticipated but not specified

**AI Analysis:**
This change adds a feature that was mentioned in the "Open Questions" section 
of the specification. The implementation is well-designed and follows security 
best practices for TOTP-based 2FA. The specification should be updated to 
formally document this feature.

---

## Decision

**Action:** Approved  
**Approved By:** ssdajoker  
**Date:** 2026-01-07  
**Reasoning:** Feature was anticipated in spec, implementation is solid, and 
it enhances security. Specification has been updated to reflect this addition.

---

## Specification Updates

The following sections of `user_authentication.md` were updated:

### 2.1 Functional Requirements

Added:
```markdown
6. **Two-Factor Authentication (2FA)**
   - TOTP-based 2FA (RFC 6238)
   - QR code generation for easy setup
   - Backup codes for account recovery (future)
   - Optional enforcement per user
   - Compatible with Google Authenticator, Authy, etc.
```

### 3.3 Data Models

Updated User model:
```go
type User struct {
    ID               uuid.UUID
    Email            string
    PasswordHash     string
    EmailVerified    bool
    TwoFactorSecret  string  // TOTP secret (encrypted)
    TwoFactorEnabled bool    // Whether 2FA is enabled
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

### 3.4 API Endpoints

Added:
```markdown
**POST /auth/2fa/enable**
- Headers: `Authorization: Bearer <access_token>`
- Response: `{ "secret": "...", "qr_code": "otpauth://..." }`

**POST /auth/2fa/verify**
- Headers: `Authorization: Bearer <access_token>`
- Body: `{ "code": "123456" }`
- Response: `{ "success": true }`

**POST /auth/2fa/disable**
- Headers: `Authorization: Bearer <access_token>`
- Body: `{ "code": "123456" }`
- Response: `{ "success": true }`

**POST /auth/login** (updated)
- Body: `{ "email": "...", "password": "...", "totp_code": "123456" }`
- Note: `totp_code` required if 2FA is enabled
```

---

## Related Issues

- None (change was approved)

---

## Notes

- Consider adding backup codes in future iteration
- Consider adding SMS-based 2FA as alternative
- Monitor adoption rate and user feedback

---

**Change Order Generated by Factory v1.0.0**
```

---

## 4. RESCUE Mode

### 4.1 Purpose

Reverse-engineer an existing codebase to generate specifications, useful for 
projects that have code but no documentation, or for understanding legacy systems.

### 4.2 User Journey

```
Start RESCUE → Scan Codebase → Infer Architecture → Generate Specs → Review & Edit → Save
```

### 4.3 Interface Design

**Scanning Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         RESCUE MODE                                  ║
║                  Reverse-Engineer Codebase                           ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 🔍 Scanning codebase...                                             │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ Repository: /home/user/projects/my-app                               │
│ Language: Go (detected)                                              │
│                                                                      │
│ Progress: [████████████████████] 100% (234/234 files)               │
│                                                                      │
│ Discovered:                                                          │
│  • 234 source files                                                  │
│  • 45 packages                                                       │
│  • 12 API endpoints                                                  │
│  • 8 database tables                                                 │
│  • 156 functions                                                     │
│  • 23 external dependencies                                          │
│                                                                      │
│ Analyzing architecture... ⏳                                         │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

**Architecture Inference Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         RESCUE MODE                                  ║
║                  Inferring System Architecture                       ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 🤖 AI is analyzing your codebase...                                 │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ ✓ Step 1/6: Analyzing project structure                             │
│ ✓ Step 2/6: Identifying components and modules                      │
│ ✓ Step 3/6: Mapping dependencies                                    │
│ ⏳ Step 4/6: Inferring data models...                               │
│ ⏸️  Step 5/6: Documenting API endpoints                             │
│ ⏸️  Step 6/6: Generating specifications                             │
│                                                                      │
│ Estimated time remaining: 2 minutes                                  │
│                                                                      │
│ Current Focus: Analyzing database schema and ORM models...           │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

**Generated Specs Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         RESCUE MODE                                  ║
║                  Generated Specifications                            ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 📚 Generated 5 specification documents:                             │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ ✓ system_architecture.md (3.2 KB)                                   │
│   High-level system design, components, and interactions             │
│                                                                      │
│ ✓ api_endpoints.md (5.8 KB)                                         │
│   Complete API documentation with 12 endpoints                       │
│                                                                      │
│ ✓ database_schema.md (4.1 KB)                                       │
│   Data models, relationships, and migrations                         │
│                                                                      │
│ ✓ authentication_system.md (6.3 KB)                                 │
│   User authentication and authorization                              │
│                                                                      │
│ ✓ dependencies.md (2.1 KB)                                          │
│   External libraries and services                                    │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ 📊 Code Quality Analysis:                                           │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ • Test Coverage: 67% (good)                                          │
│ • Code Complexity: Medium                                            │
│ • Technical Debt: Low-Medium                                         │
│ • Security Issues: 2 potential issues found                          │
│ • Performance Concerns: 1 identified                                 │
│                                                                      │
│ [View Details] [Edit Specs] [Save All] [Cancel]                     │
└──────────────────────────────────────────────────────────────────────┘
```

**Spec Preview Screen:**
```
╔══════════════════════════════════════════════════════════════════════╗
║                         RESCUE MODE                                  ║
║                  Preview: system_architecture.md                     ║
╚══════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────┐
│ 📄 Generated Specification                                          │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│ # System Architecture                                                │
│                                                                      │
│ ## Overview                                                          │
│ This is a Go-based web application following a layered architecture  │
│ pattern. The system consists of an HTTP API server, PostgreSQL       │
│ database, and Redis cache.                                           │
│                                                                      │
│ ## Architecture Pattern                                              │
│ **Pattern:** Layered Architecture (3-tier)                           │
│ - Presentation Layer: HTTP handlers                                  │
│ - Business Logic Layer: Services                                     │
│ - Data Access Layer: Repositories                                    │
│                                                                      │
│ ## Components                                                        │
│                                                                      │
│ ### 1. API Server (`cmd/server`)                                     │
│ Main entry point for the application. Initializes HTTP server,       │
│ database connections, and routes.                                    │
│                                                                      │
│ ### 2. Authentication Module (`internal/auth`)                       │
│ Handles user authentication, JWT token management, and OAuth         │
│ integration.                                                         │
│                                                                      │
│ ### 3. User Module (`internal/user`)                                 │
│ User management, profile updates, and user-related operations.       │
│                                                                      │
│ ...                                                                  │
│                                                                      │
│ [↑↓ to scroll | E to edit | N for next spec | S to save | Q to quit]│
└──────────────────────────────────────────────────────────────────────┘
```

### 4.4 Workflow Steps

#### Step 1: Codebase Scanning
- Recursively scan project directory
- Identify programming language(s)
- Count files, packages, functions
- Detect frameworks and libraries
- Find configuration files

#### Step 2: Architecture Inference
- Analyze directory structure
- Identify architectural patterns (MVC, layered, microservices, etc.)
- Map component dependencies
- Identify entry points and main flows
- Detect design patterns

#### Step 3: Component Analysis
- For each major component:
  - Identify purpose and responsibilities
  - List public interfaces
  - Find dependencies
  - Detect configuration options

#### Step 4: Data Model Extraction
- Parse database migrations or ORM models
- Generate entity-relationship diagrams
- Document data types and constraints
- Identify relationships

#### Step 5: API Documentation
- Find HTTP handlers/controllers
- Extract routes and methods
- Infer request/response formats
- Document authentication requirements

#### Step 6: Specification Generation
- Generate markdown documents for each area
- Include code examples from codebase
- Add diagrams where helpful
- Highlight technical debt and issues

#### Step 7: Review & Edit
- Present generated specs to user
- Allow editing before saving
- Validate completeness
- Save to `/contracts/`

### 4.5 LLM Prompt Template

```markdown
You are a software architect and technical writer. Analyze the following 
codebase and generate comprehensive specification documents.

## Codebase Information:
- **Language:** {language}
- **Framework:** {framework}
- **Files:** {file_count}
- **Packages:** {package_count}
- **Lines of Code:** {loc}

## Directory Structure:
{directory_tree}

## Key Files:
{key_files_content}

## Dependencies:
{dependencies}

## Analysis Instructions:

1. **Infer Architecture:**
   - Identify architectural pattern (layered, microservices, MVC, etc.)
   - Describe high-level system design
   - Map component interactions

2. **Document Components:**
   - List major components/modules
   - Describe purpose and responsibilities
   - Identify public interfaces

3. **Extract Data Models:**
   - Document database schema
   - Describe entity relationships
   - Note data types and constraints

4. **Document APIs:**
   - List all HTTP endpoints
   - Describe request/response formats
   - Note authentication requirements

5. **Identify Technical Debt:**
   - Security vulnerabilities
   - Performance bottlenecks
   - Code smells
   - Missing tests

## Output Format:

Generate the following specification documents:

### 1. system_architecture.md
```markdown
# System Architecture

## Overview
[High-level description]

## Architecture Pattern
[Pattern name and description]

## Components
[List and describe components]

## Component Interactions
[Describe how components interact]

## Technology Stack
[List technologies used]

## Deployment Architecture
[Describe deployment setup]
```

### 2. api_endpoints.md
```markdown
# API Endpoints

## Overview
[API description]

## Authentication
[Auth mechanism]

## Endpoints

### GET /endpoint
- **Description:** [what it does]
- **Auth:** [required/optional]
- **Request:** [parameters]
- **Response:** [format]
- **Example:** [code example]
```

### 3. database_schema.md
```markdown
# Database Schema

## Overview
[Database description]

## Tables

### table_name
- **Description:** [purpose]
- **Columns:**
  - `column_name` (type) - description
- **Indexes:** [list indexes]
- **Relationships:** [foreign keys]
```

### 4. [feature]_system.md (for each major feature)
```markdown
# [Feature] System

## Overview
[Feature description]

## Requirements
[Inferred requirements]

## Implementation
[How it's implemented]

## Data Models
[Related data models]

## API Endpoints
[Related endpoints]

## Security Considerations
[Security notes]
```

### 5. dependencies.md
```markdown
# Dependencies

## External Libraries
[List with versions and purposes]

## External Services
[APIs, databases, etc.]

## Development Dependencies
[Build tools, testing frameworks]
```

### 6. technical_debt.md
```markdown
# Technical Debt Analysis

## Security Issues
[List security concerns]

## Performance Issues
[List performance concerns]

## Code Quality Issues
[List code smells]

## Missing Features
[Features that seem incomplete]

## Recommendations
[Prioritized list of improvements]
```

---

Begin analysis and generate specifications:
```

### 4.6 Generated Spec Example

```markdown
# System Architecture

**Generated:** 2026-01-07  
**Source:** /home/user/projects/my-app  
**Language:** Go  
**Framework:** Standard library + Chi router

---

## Overview

This is a Go-based web application implementing a RESTful API for user 
authentication and management. The system follows a layered architecture 
pattern with clear separation between HTTP handlers, business logic, and 
data access layers.

**Key Characteristics:**
- Monolithic architecture
- PostgreSQL for persistent storage
- Redis for caching and rate limiting
- JWT-based authentication
- RESTful API design

---

## Architecture Pattern

**Pattern:** Layered Architecture (3-tier)

```
┌─────────────────────────────────────────────────────────┐
│                  Presentation Layer                     │
│              (HTTP Handlers / Controllers)              │
│                  internal/handlers/                     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Business Logic Layer                   │
│                     (Services)                          │
│                  internal/services/                     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Data Access Layer                      │
│                   (Repositories)                        │
│                internal/repositories/                   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                      Database                           │
│                   (PostgreSQL)                          │
└─────────────────────────────────────────────────────────┘
```

---

## Components

### 1. API Server (`cmd/server`)

**Purpose:** Main entry point for the application

**Responsibilities:**
- Initialize HTTP server
- Set up database connections
- Configure middleware
- Register routes
- Handle graceful shutdown

**Key Files:**
- `cmd/server/main.go` - Entry point
- `cmd/server/routes.go` - Route definitions

**Dependencies:**
- Chi router for HTTP routing
- PostgreSQL driver
- Redis client

---

### 2. Authentication Module (`internal/auth`)

**Purpose:** Handle user authentication and authorization

**Responsibilities:**
- User registration and login
- JWT token generation and validation
- OAuth integration (Google, GitHub)
- Password hashing and verification
- Session management

**Key Files:**
- `internal/auth/service.go` - Auth business logic
- `internal/auth/handlers.go` - HTTP handlers
- `internal/auth/jwt.go` - JWT utilities
- `internal/auth/oauth.go` - OAuth providers

**Public Interface:**
```go
type Service interface {
    Register(ctx context.Context, email, password string) (*User, error)
    Login(ctx context.Context, email, password string) (string, error)
    ValidateToken(ctx context.Context, token string) (*Claims, error)
    RefreshToken(ctx context.Context, refreshToken string) (string, error)
}
```

---

### 3. User Module (`internal/user`)

**Purpose:** User profile management

**Responsibilities:**
- Get user profile
- Update user information
- Delete user account
- List users (admin)

**Key Files:**
- `internal/user/service.go` - User business logic
- `internal/user/handlers.go` - HTTP handlers
- `internal/user/repository.go` - Database access

**Public Interface:**
```go
type Service interface {
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    Update(ctx context.Context, id uuid.UUID, updates *UpdateRequest) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filters *ListFilters) ([]*User, error)
}
```

---

### 4. Database Module (`internal/database`)

**Purpose:** Database connection and migration management

**Responsibilities:**
- Initialize database connection
- Run migrations
- Provide connection pool
- Handle transactions

**Key Files:**
- `internal/database/db.go` - Connection setup
- `internal/database/migrations/` - SQL migrations

---

### 5. Middleware (`internal/middleware`)

**Purpose:** HTTP middleware for cross-cutting concerns

**Responsibilities:**
- Authentication verification
- Request logging
- CORS handling
- Rate limiting
- Error recovery

**Key Files:**
- `internal/middleware/auth.go` - Auth middleware
- `internal/middleware/logging.go` - Request logging
- `internal/middleware/ratelimit.go` - Rate limiting

---

## Component Interactions

### User Registration Flow

```
Client → POST /auth/register
    ↓
AuthHandler.Register()
    ↓
AuthService.Register()
    ├─→ Validate input
    ├─→ Hash password
    └─→ UserRepository.Create()
            ↓
        Database (INSERT)
            ↓
        Return User
    ↓
Generate JWT token
    ↓
Return token to client
```

### Authenticated Request Flow

```
Client → GET /users/me (with JWT)
    ↓
AuthMiddleware.Verify()
    ├─→ Extract token from header
    ├─→ Validate JWT signature
    └─→ Extract user ID from claims
    ↓
UserHandler.GetProfile()
    ↓
UserService.GetByID()
    ↓
UserRepository.FindByID()
    ↓
Database (SELECT)
    ↓
Return user data to client
```

---

## Technology Stack

### Backend
- **Language:** Go 1.21
- **HTTP Router:** Chi v5
- **Database:** PostgreSQL 14
- **Cache:** Redis 6
- **Authentication:** JWT (golang-jwt/jwt)
- **Password Hashing:** bcrypt

### External Services
- **OAuth Providers:** Google, GitHub
- **Email:** SMTP (configurable)

### Development Tools
- **Testing:** Go testing package + testify
- **Migrations:** golang-migrate
- **Linting:** golangci-lint

---

## Deployment Architecture

**Current Setup:** Single server deployment

```
┌─────────────────────────────────────────────────────────┐
│                      Load Balancer                      │
│                       (Nginx)                           │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Application Server                   │
│                      (Go Binary)                        │
│                    Port: 8080                           │
└─────────────────────────────────────────────────────────┘
                          │
            ┌─────────────┴─────────────┐
            ▼                           ▼
┌───────────────────────┐   ┌───────────────────────┐
│     PostgreSQL        │   │        Redis          │
│     Port: 5432        │   │      Port: 6379       │
└───────────────────────┘   └───────────────────────┘
```

**Scalability Considerations:**
- Stateless design allows horizontal scaling
- Database connection pooling (max 100 connections)
- Redis for distributed rate limiting
- JWT tokens eliminate need for session storage

---

## Configuration

**Environment Variables:**
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `JWT_SECRET` - Secret for JWT signing
- `OAUTH_GOOGLE_CLIENT_ID` - Google OAuth client ID
- `OAUTH_GOOGLE_CLIENT_SECRET` - Google OAuth secret
- `OAUTH_GITHUB_CLIENT_ID` - GitHub OAuth client ID
- `OAUTH_GITHUB_CLIENT_SECRET` - GitHub OAuth secret
- `SMTP_HOST` - SMTP server host
- `SMTP_PORT` - SMTP server port
- `SMTP_USER` - SMTP username
- `SMTP_PASS` - SMTP password

---

## Security Considerations

1. **Authentication:**
   - JWT tokens with 15-minute expiry
   - Refresh tokens with 7-day expiry
   - Secure password hashing (bcrypt cost 10)

2. **Authorization:**
   - Role-based access control (RBAC)
   - Middleware enforces authentication

3. **Data Protection:**
   - HTTPS only (enforced by Nginx)
   - Secure cookie flags
   - CORS configured

4. **Rate Limiting:**
   - Implemented on login endpoint (5 per 15 min)
   - Redis-backed for distributed systems

---

## Performance Characteristics

**Observed Performance:**
- Login endpoint: ~150ms (p95)
- User profile fetch: ~50ms (p95)
- Database connection pool: 50 connections
- Concurrent users: Tested up to 500

**Bottlenecks:**
- Database queries (no caching for user profiles)
- Password hashing on login (bcrypt is intentionally slow)

---

## Technical Debt

See `technical_debt.md` for detailed analysis.

**Summary:**
- Missing rate limiting on some endpoints
- No caching for frequently accessed data
- Test coverage at 67% (target: 80%)
- Some error handling could be improved

---

**Generated by Factory v1.0.0 RESCUE Mode**
```

---

## 5. Cross-Mode Features

### 5.1 Spec Templates

**Purpose:** Provide starting points for common specification types

**Templates:**
- `feature_spec.md` - General feature specification
- `api_spec.md` - API endpoint documentation
- `database_spec.md` - Database schema specification
- `architecture_spec.md` - System architecture document
- `security_spec.md` - Security requirements
- `performance_spec.md` - Performance requirements

### 5.2 Spec Validation

**Purpose:** Ensure specifications are complete and well-formed

**Checks:**
- Required sections present
- Consistent formatting
- Valid markdown syntax
- Cross-references resolved
- Examples are valid

### 5.3 Spec Versioning

**Purpose:** Track specification evolution over time

**Features:**
- Semantic versioning (1.0.0, 1.1.0, 2.0.0)
- Changelog generation
- Diff between versions
- Rollback to previous version

### 5.4 Collaboration Features

**Purpose:** Enable team collaboration on specifications

**Features:**
- Comments and annotations
- Review workflow (draft → review → approved)
- Approval tracking
- Notification on changes

---

**End of Specification**
