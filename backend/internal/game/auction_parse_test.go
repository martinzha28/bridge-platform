package game

import "testing"

func TestParseCallBids(t *testing.T) {
	tests := []struct {
		input string
		want  Call
	}{
		{"1C", BidCall(1, ClubStrain)},
		{"1D", BidCall(1, DiamondStrain)},
		{"1H", BidCall(1, HeartStrain)},
		{"1S", BidCall(1, SpadeStrain)},
		{"1NT", BidCall(1, NoTrump)},
		{"3nt", BidCall(3, NoTrump)},
		{"7S", BidCall(7, SpadeStrain)},
	}

	for _, tt := range tests {
		call, ok := ParseCall(tt.input)
		if !ok {
			t.Errorf("ParseCall(%q) returned !ok", tt.input)
			continue
		}
		if call != tt.want {
			t.Errorf("ParseCall(%q) = %v, want %v", tt.input, call, tt.want)
		}
	}
}

func TestParseCallSpecials(t *testing.T) {
	tests := []struct {
		input    string
		wantType CallType
	}{
		{"P", Pass},
		{"PASS", Pass},
		{"p", Pass},
		{"X", Double},
		{"DBL", Double},
		{"XX", Redouble},
		{"RDBL", Redouble},
	}

	for _, tt := range tests {
		call, ok := ParseCall(tt.input)
		if !ok {
			t.Errorf("ParseCall(%q) returned !ok", tt.input)
			continue
		}
		if call.Type != tt.wantType {
			t.Errorf("ParseCall(%q).Type = %v, want %v", tt.input, call.Type, tt.wantType)
		}
	}
}

func TestParseCallInvalid(t *testing.T) {
	invalids := []string{"", "Z", "8C", "0S", "1Z", "ABC"}
	for _, s := range invalids {
		if _, ok := ParseCall(s); ok {
			t.Errorf("ParseCall(%q) should fail", s)
		}
	}
}
