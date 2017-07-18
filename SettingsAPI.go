package arn

import (
	"reflect"

	"github.com/aerogo/aero"
)

// Authorize returns an error if the given API POST request is not authorized.
func (settings *Settings) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// Update updates the settings object.
func (settings *Settings) Update(ctx *aero.Context, data interface{}) error {
	updates := data.(map[string]interface{})
	return SetObjectProperties(settings, updates, func(fullKeyName string, field *reflect.StructField, property *reflect.Value, newValue reflect.Value) (bool, error) {
		switch fullKeyName {
		case "Avatar.Source":
			settings.Avatar.Source = newValue.String()
			settings.Save() // Save needed here because RefreshAvatar fetches the settings on another server
			settings.User().RefreshAvatar()
			return true, nil

		case "Avatar.SourceURL":
			settings.Avatar.SourceURL = newValue.String()
			settings.Save() // Save needed here because RefreshAvatar fetches the settings on another server
			settings.User().RefreshAvatar()
			return true, nil

		default:
			return false, nil
		}
	})
}

// Save saves the settings in the database.
func (settings *Settings) Save() error {
	return DB.Set("Settings", settings.UserID, settings)
}