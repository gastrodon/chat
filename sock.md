# Outgoing

### Ping the websocket server
| field         | description
| ---           | ---
| type          | `ping`
| time          | ping timestamp. `Pong` from the server should be the same timestamp

### Pong a server ping
| field         | description
| ---           | ---
| type          | `pong`
| time          | Server timestamp received from this ping

### Send a message to a room

| field         | description
| ---           | ---  
| type          | `message`
| key           | Session key of this user instance
| content       | String content of this message
| room_id       | Message destination room UUID

A user must have joined a room and have permission to send messages

### Read a message in a room

| field         | description
| ---           | ---
| type          | `read`
| message_id    | Message ID in it's Room

Mark a message in a room as being read

### Invite a user to a room

| field         | description
| ---           | ---
| type          | `invite`
| key           | Session key of this user instance
| user_id       | User UUID to invite
| room_id       | Room UUID to invite to

This tells a user that they should `enter` a room.
If this room is closed, only room administrators may invite users

# Incoming

### Ping the client
| field         | description
| ---           | ---
| type          | `ping`
| time          | ping timestamp. `Pong` from the client should be the same timestamp

### Pong a client ping
| field         | description
| ---           | ---
| type          | `pong`
| time          | Client timestamp received from this ping

### Receive a message in a room

| field         | description
| ---           | ---
| type          | `message`
| time          | Message send UNIX timestamp
| content       | String content of this message
| user_id       | Message author UUID
| message_id    | Message ID in it's Room
| room_id       | Incoming message's Room UUID

Received by all users listening for messages in a room.
Messages should only be stored by the client, and never the server.

### Receive a message read receipt

| field         | description
| ---           | ---
| type          | `read`
| message_id    | Message ID in it's room
| user_id       | User UUID who did read this message

Received when a
Handled by the client entirely

### Receive a room invite

| field         | description
| ---           | ---
| type          | `invite`
| user_id       | Inviter UUID
| room_id       | Room invited to UUID

User should then `enter` this room if they choose to do so. If this room is closed, it will be open to invited users.

### Be kicked from a room

| field         | description
| ---           | ---
| type          | `kick`
| room_id       | Room kicked from UUID
| message       | Message for exit, if any given
| permanent     | `true` \| `false`

It is up to the client implementation to decide how to handle kicks,
but a kicked user will not be able to listen or send to that room

### User join or exit a room

| field         | description
| ---           | ---
| type          | `enter` \| `exit`
| user_id       | User UUID who did enter or exit this room
| room_id       | Room UUID that had a user enter or exit
