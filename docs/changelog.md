# 0.0.15
- Fix spotify track duration

# 0.0.14
- Add duplicate cache checker for jellyfin/multiscrobbler

# 0.0.13
- Fix multiscrobbler support

# 0.0.12
- Add client TZ support + Selectable on user page
- Move token auth to GET ?key=XYZ for wider webhook support
- Add Multiscrobbler support
- Add /api/v1/serverinfo for version information

# 0.0.11
- Fix redirects to /login for auth required pages
- Add handling for 401/429 + No connection responses in API calls
- Add background workers for Go (clear out password resets)
- Add spotify scrobbling!!!11111!!!!!
- Fixed timezone issues

# 0.0.10
- Fixed looking up invalid profiles
- Added valid error handling to bad request && rate limiting
- Add Sendgrid library (Will add SMTP later)
- Complete password reset process

# 0.0.9
- Fix mobile menu auto collapse on select
- Add /u/ route for public user profiles (Added private flag to db - to implement later)
- Add /user route for your own profile / edit profile
- Added handling for if API is offline/incorrect
- Add index.html loading spinner while react bundle downloads
- Change HashRouter to BrowserRouter
- Added sources column to scrobbles

# 0.0.8
- Added Admin/Site config page in frontend for admin users
- Added API POST/GET /config endpointnpm install react-select-timezone

# 0.0.7
- Switch redux -> Context
- Remove excess packages

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
