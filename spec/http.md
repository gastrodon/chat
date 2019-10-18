### Crate a new user instance

POST `/user`

POST data

| field         | type          | description           | required
| ---           | ---           | ---                   | ---
| username      | string        | Username to create    | true
| password      | string        | Password to register  | true

response `200` - this user was created

| field         | type          | description
| ---           | ---           | ---
| user_id       | string        | UUID of the created user
| key           | string        | Session key

response `400` - bad or malformed request

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `403` - you may not POST a new user

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Modify a username

PUT `/user`

params

| field         | type          | description   | request
| ---           | ---           | ---           | ---
| key           | string        | Session key   | true

PUT data

| field         | type          | description           | required
| ---           | ---           | ---                   | ---
| username      | string        | New username to use   | true

response `200` - this username was modified

| field         | type          | description
| ---           | ---           | ---
| username      | string        | New username

response `400` - bad or malformed request

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `401` - you are not authorized

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description



### Get information about some user

GET `/user/{user_id}`

params

| field         | type          | description   | request
| ---           | ---           | ---           | ---
| key           | string        | Session key   | false

response `200` - user information

| field         | type          | description
| ---           | ---           | ---
| username      | string        | Username of this user
| user_id       | string        | UUID of this user

### Log in and get a session key

POST `/key`

POST data

| field         | type          | description                   | required
| ---           | ---           | ---                           | ---
| user_id       | string        | UUID of the user to log in as | true
| password      | string        | Password to log in with       | true

response `200` - user was logged in

| field         | type          | description
| ---           | ---           | ---
| key           | String        | Session key

response `400` - bad or malformed request

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `401` - you may not log in with these credentials

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Revoke and invalidate a session key

DELETE `/key`

DELETE data

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| key           | string        | session key to revoke             | true
| user_id       | string        | UUID of the user with this key  | true

### Get information about some room

GET `/room/{room_id}`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | false

response `200` - this room exists and is visible

| field         | type          | description       
| ---           | ---           | ---
| open          | bool          | Is this room joinable by any user?
| name          | string        | Name of this room
| user_count    | integer       | User count in this room
| owner         | string        | UUID of the owner of this room
| room_id       | string        | UUID of this room

response `404` - this room does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Get the users of a room

GET `/room/{room_id}/users`

params

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| key           | string        | Session key                       | false
| limit         | integer       | Number of users to fetch          | false
| offset        | integer       | Offset index to start fetching    | false

response `200` - a section of users of this room

| field         | type              | description
| ---           | ---               | ---
| users         | array\<string\>   | Array of user UUID's
| index         | integer           | Starting index of this slice
| user_count    | integer           | User count in this room

response `404` - this room does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Get a user in a room

GET `/room/{room_id}/users/{user_id}`

params

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| key           | string        | Session key                       | false

response `200` - an instance of this user in this room

| field         | type          | description
| ---           | ---           | ---
| present       | bool          | Is this user present in this room?

response `404` - user does not exist in this room

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Get the admins of a room

GET `/room/{room_id}/admins`

params

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| key           | string        | Session key                       | false
| limit         | integer       | Number of users to fetch          | false
| offset        | integer       | Offset index to start fetching    | false

response `200` - a section of admins of this room

| field         | type              | description
| ---           | ---               | ---
| users         | array\<string\>   | Array of user UUID's
| index         | integer           | Starting index of this slice
| user_count    | integer           | User count in this room

response `404` - this room does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Get the banned users of a room

GET `/room/{room_id}/bans`

params

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| key           | string        | Session key                       | false
| limit         | integer       | Number of users to fetch          | false
| offset        | integer       | Offset index to start fetching    | false

response `200` - a section of admins of this room

| field         | type              | description
| ---           | ---               | ---
| users         | array\<string\>   | Array of user UUID's
| index         | integer           | Starting index of this slice
| user_count    | integer           | User count in this room

response `404` - this room does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Get the owner of a room

GET `/room/{room_id}/owner`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | false

response `200` - the owner of this room

| field         | type          | description
| ---           | ---           | ---
| owner         | string        | Room owner UUID

response `404` - this room does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Create a room

POST `/room/`

params

| field         | type          | description                                           | required
| ---           | ---           | ---                                                   | ---
| key           | string        | Session key                                           | true


POST data

| field         | type          | description                                           | required
| ---           | ---           | ---                                                   | ---
| name          | string        | Name of this room                                     | false
| open          | bool          | Is this room joinable by any user? default is true    | false

response `200` - the room that was created

| field         | type          | description       
| ---           | ---           | ---
| open          | bool          | Is this room joinable by any user?
| name          | string        | Name of this room
| user_count    | integer       | User count in this room
| owner         | string        | UUID of the owner of this room
| room_id       | string        | UUID of this room

### Change the information of a room

PUT `/room/`

params

| field         | type          | description                                           | required
| ---           | ---           | ---                                                   | ---
| key           | string        | Session key                                           | true


POST data
At least one field is required, but not both

| field         | type          | description                                           | required
| ---           | ---           | ---                                                   | ---
| name          | string        | Name of this room                                     | false
| open          | bool          | Is this room joinable by any user? default is true    | false

response `200` - the room that was modified

| field         | type          | description
| ---           | ---           | ---
| open          | bool          | Is this room joinable by any user?
| name          | string        | Name of this room
| user_count    | integer       | User count in this room
| owner         | string        | UUID of the owner of this room
| room_id       | string        | UUID of this room

response `403` - this user may not modify this room
Note: only room owners and admins may modify a room

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `409` - this user is already an admin in this room

response `404` - this chat does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `409` - this user is already an admin in this room

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Change the owner of a room

POST `/room/{room_id}/owner`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | true

POST data

| field         | type          | description                   | required
| ---           | ---           | ---                           | ---
| user_id       | string        | UUID of the new room owner    | true

response `204` - this Room info with a new owner

response `403` - the user who did make this request may not change the room owner.
Note: only room owners may change room ownership

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `404` - this user was not found in this room, or this chat does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Add an admin to a room

PUT `/room/{room_id}/admins`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | true

PUT data

| field         | type          | description                   | required
| ---           | ---           | ---                           | ---
| user_id       | string        | UUID of the new room admin    | true

response `204` - this user is now a room admin

response `403` - this user may not add room admins
Note: only room owners may add room admins

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `404` - this user was not found in this room, or this chat does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `409` - this user is already an admin in this room

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Remove an admin from a room

DELETE `/room/{room_id}/admins`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | true

PUT data

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| user_id       | string        | UUID of the room admin to remove  | true

response `204` - this user is no longer a room admin

response `403` - this user may not remove room admins
Note: only room owners may remove room admins

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `404` - this user was not found in this room, or this chat does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `409` - this user is not an admin in this room

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Invite a user to a room

PUT `/room/{room_id}/users`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | true

PUT data

| field         | type          | description                   | required
| ---           | ---           | ---                           | ---
| user_id       | string        | UUID of the user to invite    | true

response `204` - this user was invited to this room

response `403` - this user may not invite users
Note: in non-open rooms, only admins and the room owner may invite users

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `404` - this user does not exist, or this chat does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `409` - this user is already in this room

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Remove a user from a room

DELETE `/room/{room_id}/users`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | true

PUT data

| field         | type          | description                   | required
| ---           | ---           | ---                           | ---
| user_id       | string        | UUID of the user to remove    | true

response `204` - this user was removed from this room

response `403` - this user may not remove users
Note: only admins and the room owner may remove users

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `404` - this user does not exist, or this chat does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `409` - this user is not in this room

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Ban a user from a room

PUT `/room/{room_id}/bans`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | true

PUT data

| field         | type          | description                   | required
| ---           | ---           | ---                           | ---
| user_id       | string        | UUID of the user to ban       | true

response `204` - this user was banned from this room

response `403` - this user may not ban users
Note: only admins and the room owner may ban users

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `404` - this user does not exist, or this chat does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `409` - this user is already banned in this room

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Lift a ban of a user in a room

DELETE `/room/{room_id}/bans`

params

| field         | type          | description   | required
| ---           | ---           | ---           | ---
| key           | string        | Session key   | true

PUT data

| field         | type          | description                           | required
| ---           | ---           | ---                                   | ---
| user_id       | string        | UUID of the user to lift the ban of   | true

response `204` - this ban was lifted from this room

response `403` - this user may not lift bans
Note: only admins and the room owner may lift bans

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `404` - this user does not exist, or this chat does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

response `409` - this user is not banned in this server

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description
