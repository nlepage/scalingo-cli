package users

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v7"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/gopassword"
)

func UpdateUser(ctx context.Context, app, addonUUID, username string) error {
	isSupported, err := doesDatabaseHandleUserManagement(ctx, app, addonUUID)
	if err != nil {
		return errors.Wrap(ctx, err, "get user management information")
	}

	if !isSupported {
		io.Error(ErrDatabaseNotSupportUserManagement)
		return nil
	}

	if usernameValidation, ok := isUsernameValid(username); !ok {
		io.Error(usernameValidation)
		return nil
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	l, err := c.DatabaseListUsers(ctx, app, addonUUID)
	if err != nil {
		return errors.Wrap(ctx, err, "list the database's users")
	}
	// Check if the userParam exists and is not protected
	var user *scalingo.DatabaseUser
	for _, u := range l {
		if u.Name == username {
			user = &u
			break
		}
	}
	if user == nil {
		return errors.New(ctx, fmt.Sprintf("User \"%s\" does not exist", username))
	}
	if user.Protected {
		return errors.New(ctx, fmt.Sprintf("User \"%s\" is protected", username))
	}

	password, confirmedPassword, err := askForPasswordWithRetry(ctx, 3)
	if err != nil {
		io.Error(err)
		return nil
	}

	isPasswordGenerated := false
	if password == "" {
		isPasswordGenerated = true
		password = gopassword.Generate(64)
		confirmedPassword = password
	}

	userUpdateParam := scalingo.DatabaseUpdateUserParam{
		DatabaseID:           addonUUID,
		Password:             password,
		PasswordConfirmation: confirmedPassword,
	}
	databaseUsers, err := c.DatabaseUpdateUser(ctx, app, addonUUID, username, userUpdateParam)
	if err != nil {
		return errors.Wrap(ctx, err, "update password of the given database user")
	}

	if isPasswordGenerated {
		fmt.Printf("User \"%s\" updated with password \"%s\".\n", databaseUsers.Name, password)
		return nil
	}

	fmt.Printf("User \"%s\" password updated.\n", databaseUsers.Name)
	return nil
}
