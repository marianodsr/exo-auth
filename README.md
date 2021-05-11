# exo-auth

Auth microservice for EXO app.

It handles user crud, registering and signin in.

Creates an access JWT and a refresh one when signing in.

The refresh token is placed in an http only cookie and the acess token is sent to the client and expected in an authentication header
in further requests.

Checks for authenticated users and refresh their access tokens if needed.
