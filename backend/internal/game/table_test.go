package game

import "testing"

func TestDirectionNext(t *testing.T) {
	if North.Next() != East {
		t.Errorf("North.Next() = %v, want East", North.Next())
	}
	if East.Next() != South {
		t.Errorf("East.Next() = %v, want South", East.Next())
	}
	if South.Next() != West {
		t.Errorf("South.Next() = %v, want West", South.Next())
	}
	if West.Next() != North {
		t.Errorf("West.Next() = %v, want North", West.Next())
	}
}

func TestDirectionPartner(t *testing.T) {
	if North.Partner() != South {
		t.Errorf("North.Partner() = %v, want South", North.Partner())
	}
	if East.Partner() != West {
		t.Errorf("East.Partner() = %v, want West", East.Partner())
	}
}

func TestVulnerability(t *testing.T) {
	if VulNone.IsVulnerable(North) {
		t.Error("VulNone should not make North vulnerable")
	}
	if !VulNS.IsVulnerable(North) {
		t.Error("VulNS should make North vulnerable")
	}
	if !VulNS.IsVulnerable(South) {
		t.Error("VulNS should make South vulnerable")
	}
	if VulNS.IsVulnerable(East) {
		t.Error("VulNS should not make East vulnerable")
	}
	if !VulEW.IsVulnerable(East) {
		t.Error("VulEW should make East vulnerable")
	}
	if !VulBoth.IsVulnerable(West) {
		t.Error("VulBoth should make West vulnerable")
	}
}
