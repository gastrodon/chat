[![Travis (.org)](https://img.shields.io/travis/basswaver/chat?logo=travis&style=flat-square)](https://travis-ci.org/basswaver/chat)

## Chat
A small chat. Or rather a design for a small chat.

### About
I do not like the insecure design of chats that require a lot of information to use, and that store my messages. So I started to play around with my own idea for a chat API design

### Design Goals
##### Ephemeral
Messages should only exist on the device of the sender and intended recipient. The server should just be a relay that keeps track of where to send messages, and should not remember what is sent, and ideally should be able to destroy that information as soon as it's no longer in use.

##### Minimal
As mentioned, the chat should not do much more than keep track of who is sending messages where. In addition, since there are rooms (where multiple people speak) this implementation will keep track of who can modify and remove users from rooms. A room consists of:

- A ID
- An owner
- A list of members
- A list of admins
- A list of bans

A user consists of:

- A ID
- A non-unique username (optional)

Because you probably don't need much else. The purpose of having a ID bound to a user is so that a room can keep track of who's in it, which tells the server which socket connections to send messages to. Only closed rooms might make use of bans, as open rooms are easily reentered with new accounts by banned users. So close your doors!

##### Anonymous
A user is represented only by an ID in a room, which in turn is just an ID that keeps track of some more ID's. A user can have a username, or they can just be `Anonymous`

##### Small
This repo has an implementation of this chat server, but the important part is the spec. It should be small enough to be implemented in a little while, so that it can be amended to fit the needs of the user. There shouldn't be time wasted building things that are non-essential
