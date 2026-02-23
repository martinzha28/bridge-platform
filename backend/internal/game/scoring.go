package game

const (
	minorTrickValue = 20
	majorTrickValue = 30
)

func Score(contract Contract, vulnerable bool, tricksTaken int) int {
	target := int(contract.Level) + 6

	if tricksTaken >= target {
		return makeScore(contract, vulnerable, tricksTaken-target)
	}
	return undertrickScore(contract, vulnerable, target-tricksTaken)
}

func makeScore(contract Contract, vulnerable bool, overtricks int) int {
	score := 0

	perTrick := majorTrickValue
	if contract.Strain == ClubStrain || contract.Strain == DiamondStrain {
		perTrick = minorTrickValue
	}

	trickScore := perTrick * int(contract.Level)
	if contract.Strain == NoTrump {
		trickScore += 10 // NT gets +10 for the first trick
	}

	if contract.Redoubled {
		trickScore *= 4
	} else if contract.Doubled {
		trickScore *= 2
	}
	score += trickScore

	// Game bonus
	if trickScore >= 100 {
		if vulnerable {
			score += 500
		} else {
			score += 300
		}
	} else {
		score += 50 // partscore
	}

	// Slam bonuses
	switch contract.Level {
	case 6:
		if vulnerable {
			score += 750
		} else {
			score += 500
		}
	case 7:
		if vulnerable {
			score += 1500
		} else {
			score += 1000
		}
	}

	// Overtrick score
	if contract.Redoubled {
		if vulnerable {
			score += overtricks * 400
		} else {
			score += overtricks * 200
		}
	} else if contract.Doubled {
		if vulnerable {
			score += overtricks * 200
		} else {
			score += overtricks * 100
		}
	} else {
		score += overtricks * perTrick
	}

	// Insult bonus for making doubled/redoubled
	if contract.Redoubled {
		score += 100
	} else if contract.Doubled {
		score += 50
	}

	return score
}

func undertrickScore(contract Contract, vulnerable bool, down int) int {
	score := 0

	if contract.Redoubled {
		for i := 1; i <= down; i++ {
			if vulnerable {
				if i == 1 {
					score += 400
				} else {
					score += 600
				}
			} else {
				if i == 1 {
					score += 200
				} else if i <= 3 {
					score += 400
				} else {
					score += 600
				}
			}
		}
	} else if contract.Doubled {
		for i := 1; i <= down; i++ {
			if vulnerable {
				if i == 1 {
					score += 200
				} else {
					score += 300
				}
			} else {
				if i == 1 {
					score += 100
				} else if i <= 3 {
					score += 200
				} else {
					score += 300
				}
			}
		}
	} else {
		if vulnerable {
			score += down * 100
		} else {
			score += down * 50
		}
	}

	return -score
}
