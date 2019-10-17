package models

type Room struct {
    Name string
    OwnerId string
    ID string
    Open bool
    UserCount uint16
    UserArray []string
    AdminArray []string
    InviteArray []string
}
