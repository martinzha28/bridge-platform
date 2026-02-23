package game

import "testing"

func TestScorePartscoreMade(t *testing.T) {
	// NV: 2H = -> 110
	contract := Contract{Level: 2, Strain: HeartStrain}
	score := Score(contract, false, 8)
	if score != 110 {
		t.Errorf("2H= NV = %d, want 110", score)
	}
}

func TestScoreGameMade(t *testing.T) {
	// NV: 3NT = -> 400
	contract := Contract{Level: 3, Strain: NoTrump}
	score := Score(contract, false, 9)
	if score != 400 {
		t.Errorf("3NT= NV = %d, want 400", score)
	}
}

func TestScoreGameVulnerable(t *testing.T) {
	// V: 4H = -> 620
	contract := Contract{Level: 4, Strain: HeartStrain}
	score := Score(contract, true, 10)
	if score != 620 {
		t.Errorf("4H= V = %d, want 620", score)
	}
}

func TestScoreMinorGameMade(t *testing.T) {
	// NV: 5C = -> 400
	contract := Contract{Level: 5, Strain: ClubStrain}
	score := Score(contract, false, 11)
	if score != 400 {
		t.Errorf("5C= NV = %d, want 400", score)
	}
}

func TestScoreOvertricks(t *testing.T) {
	// NV: 3NT + 1 -> 430
	contract := Contract{Level: 3, Strain: NoTrump}
	score := Score(contract, false, 10)
	if score != 430 {
		t.Errorf("3NT+1 NV = %d, want 430", score)
	}
}

func TestScoreMinorOvertrick(t *testing.T) {
	// NV: 2D + 1 -> 110
	contract := Contract{Level: 2, Strain: DiamondStrain}
	score := Score(contract, false, 9)
	if score != 110 {
		t.Errorf("2D+1 NV = %d, want 110", score)
	}
}

func TestScoreSmallSlamNV(t *testing.T) {
	// NV: 6S = -> 980
	contract := Contract{Level: 6, Strain: SpadeStrain}
	score := Score(contract, false, 12)
	if score != 980 {
		t.Errorf("6S= NV = %d, want 980", score)
	}
}

func TestScoreSmallSlamVul(t *testing.T) {
	// V: 6H = -> 1430
	contract := Contract{Level: 6, Strain: HeartStrain}
	score := Score(contract, true, 12)
	if score != 1430 {
		t.Errorf("6H= V = %d, want 1430", score)
	}
}

func TestScoreGrandSlamNV(t *testing.T) {
	// NV: 7NT = -> 1520
	contract := Contract{Level: 7, Strain: NoTrump}
	score := Score(contract, false, 13)
	if score != 1520 {
		t.Errorf("7NT= NV = %d, want 1520", score)
	}
}

func TestScoreGrandSlamVul(t *testing.T) {
	// V: 7NT = -> 2220
	contract := Contract{Level: 7, Strain: NoTrump}
	score := Score(contract, true, 13)
	if score != 2220 {
		t.Errorf("7NT= V = %d, want 2220", score)
	}
}

func TestScoreDown1NV(t *testing.T) {
	// NV: 4S -1 -> -50
	contract := Contract{Level: 4, Strain: SpadeStrain}
	score := Score(contract, false, 9)
	if score != -50 {
		t.Errorf("4S-1 NV = %d, want -50", score)
	}
}

func TestScoreDown1Vul(t *testing.T) {
	// V: 4S -1 -> -100
	contract := Contract{Level: 4, Strain: SpadeStrain}
	score := Score(contract, true, 9)
	if score != -100 {
		t.Errorf("4S-1 V = %d, want -100", score)
	}
}

func TestScoreDown3NV(t *testing.T) {
	// NV: 4S -3 -> -150
	contract := Contract{Level: 4, Strain: SpadeStrain}
	score := Score(contract, false, 7)
	if score != -150 {
		t.Errorf("4S-3 NV = %d, want -150", score)
	}
}

func TestScoreDown3Vul(t *testing.T) {
	// V: 4S -3 -> -300
	contract := Contract{Level: 4, Strain: SpadeStrain}
	score := Score(contract, true, 7)
	if score != -300 {
		t.Errorf("4S-3 V = %d, want -300", score)
	}
}

func TestScoreDoubledMade(t *testing.T) {
	// NV: 2SX = -> 470
	contract := Contract{Level: 2, Strain: SpadeStrain, Doubled: true}
	score := Score(contract, false, 8)
	if score != 470 {
		t.Errorf("2SX= NV = %d, want 470", score)
	}
}

func TestScoreDoubledOvertrick(t *testing.T) {
	// NV: 2SX +1 -> 570
	contract := Contract{Level: 2, Strain: SpadeStrain, Doubled: true}
	score := Score(contract, false, 9)
	if score != 570 {
		t.Errorf("2SX+1 NV = %d, want 570", score)
	}
}

func TestScoreDoubledOvertrickVul(t *testing.T) {
	// V: 2SX +1 -> 870
	contract := Contract{Level: 2, Strain: SpadeStrain, Doubled: true}
	score := Score(contract, true, 9)
	if score != 870 {
		t.Errorf("2SX+1 V = %d, want 870", score)
	}
}

func TestScoreDoubledDown1NV(t *testing.T) {
	// NV: 4SX -1 -> -100
	contract := Contract{Level: 4, Strain: SpadeStrain, Doubled: true}
	score := Score(contract, false, 9)
	if score != -100 {
		t.Errorf("4SX-1 NV = %d, want -100", score)
	}
}

func TestScoreDoubledDown1Vul(t *testing.T) {
	// V: 4SX -1 -> -200
	contract := Contract{Level: 4, Strain: SpadeStrain, Doubled: true}
	score := Score(contract, true, 9)
	if score != -200 {
		t.Errorf("4SX-1 V = %d, want -200", score)
	}
}

func TestScoreDoubledDown3NV(t *testing.T) {
	// NV: 4SX -3 -> -500
	contract := Contract{Level: 4, Strain: SpadeStrain, Doubled: true}
	score := Score(contract, false, 7)
	if score != -500 {
		t.Errorf("4SX-3 NV = %d, want -500", score)
	}
}

func TestScoreDoubledDown3Vul(t *testing.T) {
	// V: 4SX -3 -> -800
	contract := Contract{Level: 4, Strain: SpadeStrain, Doubled: true}
	score := Score(contract, true, 7)
	if score != -800 {
		t.Errorf("4SX-3 V = %d, want -800", score)
	}
}

func TestScoreRedoubledMade(t *testing.T) {
	// NV: 2SXX = -> 640
	contract := Contract{Level: 2, Strain: SpadeStrain, Redoubled: true}
	score := Score(contract, false, 8)
	if score != 640 {
		t.Errorf("2SXX= NV = %d, want 640", score)
	}
}

func TestScoreRedoubledDown1NV(t *testing.T) {
	// NV: 4SXX -1 -> -200
	contract := Contract{Level: 4, Strain: SpadeStrain, Redoubled: true}
	score := Score(contract, false, 9)
	if score != -200 {
		t.Errorf("4SXX-1 NV = %d, want -200", score)
	}
}

func TestScoreRedoubledDown1Vul(t *testing.T) {
	// V: 4SXX -1 -> -400
	contract := Contract{Level: 4, Strain: SpadeStrain, Redoubled: true}
	score := Score(contract, true, 9)
	if score != -400 {
		t.Errorf("4SXX-1 V = %d, want -400", score)
	}
}

func TestScoreDoubledDown4NV(t *testing.T) {
	// NV: 4SX -4 -> -800
	contract := Contract{Level: 4, Strain: SpadeStrain, Doubled: true}
	score := Score(contract, false, 6)
	if score != -800 {
		t.Errorf("4SX-4 NV = %d, want -800", score)
	}
}

func TestScore1NTJustMade(t *testing.T) {
	// NV: 1NT = -> 90
	contract := Contract{Level: 1, Strain: NoTrump}
	score := Score(contract, false, 7)
	if score != 90 {
		t.Errorf("1NT= NV = %d, want 90", score)
	}
}

func TestScore1MinorJustMade(t *testing.T) {
	// NV: 1D = -> 70
	contract := Contract{Level: 1, Strain: DiamondStrain}
	score := Score(contract, false, 7)
	if score != 70 {
		t.Errorf("1D= NV = %d, want 70", score)
	}
}

func TestScore1NTDoubledJustMade(t *testing.T) {
	// NV: 1NTX = -> 180
	contract := Contract{Level: 1, Strain: NoTrump, Doubled: true}
	score := Score(contract, false, 7)
	if score != 180 {
		t.Errorf("1NTX= NV = %d, want 180", score)
	}
}
