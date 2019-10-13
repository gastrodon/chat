### Crate a new user instance

POST `/user`

POST data

| field         | type          | description           | required
| ---           | ---           | ---                   | ---
| username      | string        | Username to create    | true
| password      | string        | Password to register  | true

response `202` - this user was accepted

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

### Get information about some user

GET `/user/{user_id}`

response `200` - user information

| field         | type          | description
| ---           | ---           | ---
| username      | string        | Username of this user
| user_id       | string        | UUID of this user

### Log in and get a session key

GET `/token`

params

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

DELETE `/token`

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| key           | string        | session key to revoke             | true
| user_id       | string        | UUID of the user with this token  | true

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
| user_count    | integer       | User count in this channel
| owner         | string        | UUID of the owner of this room

response `404` - this room does not exist

| field         | type          | description
| ---           | ---           | ---
| error         | string        | Error description

### Get the users of a room

GET `/room/{room_id}/members`

params

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| key           | string        | Session key                       | false
| limit         | integer       | Number of members to fetch        | false
| offset        | integer       | Offset index to start fetching    | false

response `200` - a section of members of this room

| field         | type              | description
| ---           | ---               | ---
| users         | array\<string\>   | Array of user UUID's
| index         | integer           | Starting index of this slice
| user_count    | integer           | User count in this channel

### Get a user in a room

GET `/room/{room_id}/members/{user_id}`

params

| field         | type          | description                       | required
| ---           | ---           | ---                               | ---
| key           | string        | Session key                       | false

response `200` - an instance of this user in this room

| field         | type          | description
| ---           | ---           | ---
| present       | bool          | Is this user present in this chat?

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
| limit         | integer       | Number of members to fetch        | false
| offset        | integer       | Offset index to start fetching    | false

response `200` - a section of admins of this room

| field         | type              | description
| ---           | ---               | ---
| users         | array\<string\>   | Array of user UUID's
| index         | integer           | Starting index of this slice
| user_count    | integer           | User count in this channel


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

### Create a room

POST `/room/`

POST data

| field         | type          | description                                           | required
| ---           | ---           | ---                                                   | ---
| key           | string        | Session key                                           | true
| name          | string        | Name of this room                                     | false
| open          | bool          | Is this room joinable by any user? default is true    | false

response - `200` - the UUID to this room

| field         | type          | description
| ---           | ---           | ---
| room_id       | Room UUID     | UUID of the created room
