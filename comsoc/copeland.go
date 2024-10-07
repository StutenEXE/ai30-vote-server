package comsoc

func CopelandSWF(p Profile) (Count, error) {
	alts := make([]Alternative, len(p[0]))
	copy(alts, p[0])

	scores := make(Count)

	// On itère sur les alternatives (sauf la dernière car elle aura déjà fait tous les duels)
	for i := 0; i < len(alts)-1; i++ {
		// On itère sur les alternatives qui n'ont pas encore été comparées à i
		for j := i + 1; j < len(alts); j++ {
			bestAlts, err := MajoritySCF(keepOnlyTwoAlternatives(p, alts[i], alts[j]))
			if err != nil {
				return Count{}, err
			}
			if len(bestAlts) == 1 && bestAlts[0] == alts[i] {
				scores[alts[i]] = scores[alts[i]] + 1
				scores[alts[j]] = scores[alts[j]] - 1
			} else if len(bestAlts) == 1 && bestAlts[0] == alts[j] {
				scores[alts[i]] = scores[alts[i]] - 1
				scores[alts[j]] = scores[alts[j]] + 1
			}
		}
	}

	return scores, nil
}

func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)
	if err != nil {
		return
	}
	bestAlts = maxCount(count)
	return
}
