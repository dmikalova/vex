package engine

import "testing"

func TestPlayPermissionValidate(t *testing.T) {
	if err := (PlayPermission{}).validate(); err != nil {
		t.Errorf("ungranted permission = %v, want nil", err)
	}
	if err := (PlayPermission{
		House:  Untamed,
		Amount: 1,
	}).validate(); err != nil {
		t.Errorf("counted permission = %v, want nil", err)
	}
	if err := (PlayPermission{House: Untamed}).validate(); err == nil {
		t.Error("a granted permission with no count should be rejected")
	}
	if err := (PlayPermission{Types: CardTypesOf(Upgrade)}).validate(); err != nil {
		t.Errorf("a Types-scoped waiver needs no count, got %v", err)
	}
	if !(PlayPermission{Types: CardTypesOf(Upgrade)}).granted() {
		t.Error("a Types-scoped waiver should be granted")
	}
}
