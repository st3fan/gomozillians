// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package mozillians

import (
	"code.google.com/p/go.net/context"
	"os"
	"testing"
)

func Test_FindUsers(t *testing.T) {
	client := NewClient(os.Getenv("MOZILLIANS_APP_NAME"), os.Getenv("MOZILLIANS_APP_KEY"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	users, _, err := client.Users(ctx, UsersOptions{Email: "sarentz@mozilla.com"})
	if err != nil {
		t.Fatal("Users() failed: ", err)
	}
	if len(users) != 1 {
		t.Fatal("Did not get expected (1) number of users")
	}
}

func Test_GetUserByEmail_UnknownUser(t *testing.T) {
	client := NewClient(os.Getenv("MOZILLIANS_APP_NAME"), os.Getenv("MOZILLIANS_APP_KEY"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	users, _, err := client.Users(ctx, UsersOptions{Email: "doesnotexist@mozilla.com"})
	if err != nil {
		t.Fatal("Users() failed: ", err)
	}
	if len(users) != 0 {
		t.Fatal("Did not get expected (0) number of users")
	}
}
