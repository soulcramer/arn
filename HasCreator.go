package arn

import (
	"time"
)

// HasCreator includes user ID and date for the creation of this object.
type HasCreator struct {
	Created   string `json:"created"`
	CreatedBy string `json:"createdBy"`
}

// Creator returns the user who created this object.
func (obj *HasCreator) Creator() *User {
	user, _ := GetUser(obj.CreatedBy)
	return user
}

// CreatorID returns the ID of the user who created this object.
func (obj *HasCreator) CreatorID() string {
	return obj.CreatedBy
}

// GetCreated returns the creation time of the object.
func (obj *HasCreator) GetCreated() string {
	return obj.Created
}

// GetCreatedBy returns the ID of the user who created this object.
func (obj *HasCreator) GetCreatedBy() string {
	return obj.CreatedBy
}

// GetCreatedTime returns the creation time of the object as a time struct.
func (obj *HasCreator) GetCreatedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, obj.Created)
	return t
}
