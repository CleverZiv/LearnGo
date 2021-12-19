package design_pattern_in_go

import "testing"

func TestNewUser(t *testing.T) {
	user, err := NewUser("1", "da", WithAge(20), WithEmail("100231"))
	if err != nil {
		t.Log(err)
	}
	t.Log(user)
}
