package engine

import "testing"

func TestPlayPermissionValidate(t *testing.T) {
	if err := (PlayPermission{}).validate(); err != nil {
		t.Errorf("ungranted permission = %v, want nil", err)
	}
	if err := (PlayPermission{House: Untamed, Amount: 1}).validate(); err != nil {
		t.Errorf("counted permission = %v, want nil", err)
	}
	if err := (PlayPermission{House: Untamed}).validate(); err == nil {
		t.Error("a granted permission with no count should be rejected")
	}
}
