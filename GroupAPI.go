package arn

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ Joinable      = (*Group)(nil)
	_ Publishable   = (*Group)(nil)
	_ PostParent    = (*Group)(nil)
	_ fmt.Stringer  = (*Group)(nil)
	_ api.Newable   = (*Group)(nil)
	_ api.Editable  = (*Group)(nil)
	_ api.Deletable = (*Group)(nil)
)

// Actions
func init() {
	API.RegisterActions("Group", []*api.Action{
		// Publish
		PublishAction(),

		// Unpublish
		UnpublishAction(),

		// Join
		JoinAction(),

		// Leave
		LeaveAction(),
	})
}

// Create ...
func (group *Group) Create(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	group.ID = GenerateID("Group")
	group.Created = DateTimeUTC()
	group.CreatedBy = user.ID
	group.Edited = group.Created
	group.EditedBy = group.CreatedBy

	group.Members = []*GroupMember{
		&GroupMember{
			UserID: user.ID,
			Role:   "founder",
			Joined: group.Created,
		},
	}

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "create", "Group", group.ID, "", "", "")
	logEntry.Save()

	return group.Unpublish()
}

// Edit creates an edit log entry.
func (group *Group) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (consumed bool, err error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "Group", group.ID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// OnAppend saves a log entry.
func (group *Group) OnAppend(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayAppend", "Group", group.ID, fmt.Sprintf("%s[%d]", key, index), "", fmt.Sprint(obj))
	logEntry.Save()
}

// OnRemove saves a log entry.
func (group *Group) OnRemove(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayRemove", "Group", group.ID, fmt.Sprintf("%s[%d]", key, index), fmt.Sprint(obj), "")
	logEntry.Save()
}

// AfterEdit updates the metadata.
func (group *Group) AfterEdit(ctx *aero.Context) error {
	group.Edited = DateTimeUTC()
	group.EditedBy = GetUserFromContext(ctx).ID
	return nil
}

// Delete deletes the object from the database.
func (group *Group) Delete() error {
	if group.IsDraft {
		draftIndex := group.Creator().DraftIndex()
		draftIndex.GroupID = ""
		draftIndex.Save()
	}

	DB.Delete("Group", group.ID)
	return nil
}

// DeleteInContext deletes the amv in the given context.
func (group *Group) DeleteInContext(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "delete", "Group", group.ID, "", fmt.Sprint(group), "")
	logEntry.Save()

	return group.Delete()
}

// Authorize returns an error if the given API POST request is not authorized.
func (group *Group) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	if action == "edit" && group.CreatedBy != user.ID {
		return errors.New("Can't edit groups from other people")
	}

	return nil
}

// Save saves the group in the database.
func (group *Group) Save() {
	DB.Set("Group", group.ID, group)
}
