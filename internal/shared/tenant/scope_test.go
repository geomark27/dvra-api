package tenant

import "testing"

func TestScopeAllowsOnlyItsCompany(t *testing.T) {
	scope := ForCompany(10)
	if !scope.Valid() || !scope.Allows(10) || scope.Allows(11) {
		t.Fatalf("unexpected company scope behavior: %#v", scope)
	}
	if !GlobalScope().Allows(11) {
		t.Fatal("global scope should allow all companies")
	}
	if (Scope{}).Valid() || (Scope{}).Allows(0) {
		t.Fatal("zero scope must fail closed")
	}
	actor := ForActor(10, 7, "hiring_manager")
	if actor.UserID != 7 || actor.Role != "hiring_manager" || !actor.Allows(10) {
		t.Fatalf("unexpected actor scope: %#v", actor)
	}
}
