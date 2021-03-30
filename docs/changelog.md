# 0.0.6
- Fix hitting dashboard when logged out
- Clean up app.js

# 0.0.5
- Only allow ItemType:Audio from Jellyfin
- Fix NavBar for Mobile (Ugly hack but.. TO REWORK)
- Fixed registration page issues
- Add functionality to pull recent scrobbles to Dashboard
- Add MX record lookup validation for emails
- Add username validation for a-Z 0-9 _ and .
- Dashboard shows basic table of last 500 scrobbles.

# 0.0.4
-  Display stats on homepage

# 0.0.3
- Clean up login/redirect flow
- Add redirect when not authed on other endpoints
- Add GET /stats endpoint for overal stats

# 0.0.2
- Login flow working..
- Jellyfin scrobble working
- Returns scrobbles via API for authed users /api/v1/user/{uuid}/scrobble
- Add redis handler + funcs
- Move middleware to pass in uuid as needed

# 0.0.1
- Initial commit
- Added basic registration/login flow
