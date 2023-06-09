# 0.2.0
- Restructure back to single repo because it's a lot easier to handle

# 0.1.8
- Add list of tracks on artist/album pages

# 0.1.7
- Fix navidrome logging 2x scrobbles if second player is paused
- Fix jellyfin logging 2x scrobbles if 2 players are playing

# 0.1.6
- Allow UUID of 0 in artists/{uuid}/top to return site-wide stats
- Set images as failed if cannot find in spotify to prevent oversearching

# 0.1.51
- Fixed image resize to maintain aspect ratio

# 0.1.5
- Add image upload endpoints for artist/albums @ v1/albums/{uuid}/upload

# 0.1.4
- Fix spotify image import on scrobble for new artists/albums
- Create image resizer

# 0.1.3
- Added date ranges to profile scrobbles

# 0.1.2
- Add docker-compose file for local dev
- Implemented top listeners for artist/album endpoints to match track
- Add recent endpoint

# 0.1.1
- Cached all config values
- Updated spotify sdk package to v2
- Changed package name to gitlab.com/goscrobble/goscrobble-api to match repo
- Updated duplicate scrobble logic to never log the same song twice

# 0.1.0
- Split frontend/backend code into separate repos (https://gitlab.com/goscrobble/goscrobble-web)
- Added new ENV VARS to support unique configurations: DATA_DIRECTORY, FRONTEND_DIRECTORY, API_DOCS_DIRECTORY
- Started API documentation @ /docs (https://goscrobble.com/docs/)
- Added docker-compose file for local development!

# 0.0.33
- Add mod permission
- Add track/edit page

# 0.0.32
- Add related records into track API
- Build out track page to show links to related records
- Tidy UI *even more*
- Bump golang build to 1.16.7
- Added DevMode env var. This prevents the background workers running on local machines

# 0.0.31
- Added newlines for flamerohr
- Tidied pages

# 0.0.30
- Fix artist/album/track images

# 0.0.29
- Add image handler
- Store images locally

# 0.0.28
- Fix mobile view on user pages
- Fix favicon issue
- Return Artist UUID with scrobble list

# 0.0.27
- Navidrome works!
- Tidy up request/response structure in backend
- Tidy Settings page

# 0.0.26
- Make email required
- Add basic navidrome/subsonic connection

# 0.0.25
- Images now pull from spotify if setup!
- Show top artists/album

# 0.0.24
- Spotify will now add images on scrobble
- Renamed /api/v1/track to tracks to bypass blocklists... (uBlock :()

# 0.0.23
- Get top tracks pulling correctly :)

# 0.0.22
- Rework navbar + user pages

# 0.0.21
- Add ez deploy script
- Half implemented JWT refresh tokens, need to finish JS implementation

# 0.0.20
- Return related data on artist/album/track endpoints
- Scrobble table now links to tracks

# 0.0.19
- Tidy init/goscrobble.service
- Add routers for Artist/Album/Track endpoints + basic pages
- Move UUID generation into Go so we don't have to query the record!! Wooo!

# 0.0.18
- Add MBID/Spotify Autolinking if track exists
- Add Genre table + .go files

# 0.0.17
- Add check for registration_enabled on /register endpoint
- Made songlookup check artist name as well

# 0.0.16
- Add registration_enabled to /api/v1/serverinfo
- Add config table caching on save
- Fix redis TTL not being parsed correctly
- Move registration enabled to backend toggle
- Fixed navbar when loading /u/profile URL
- Token now shows on user page + can reset
- Added basic popup validation to disconnect/reset buttons

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
