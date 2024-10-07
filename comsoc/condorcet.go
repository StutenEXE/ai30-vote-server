package comsoc

// Permet de ne garder que les deux alternatives passées en paramètre dans le profil
func keepOnlyTwoAlternatives(p Profile, alt1 Alternative, alt2 Alternative) (newProfile Profile) {
	for i, alts := range p {
		newProfile = append(newProfile, make([]Alternative, len(alts)))
		copy(newProfile[i], alts)
		numberOfDeletions := 0
		for j, alt := range alts {
			if alt != alt1 && alt != alt2 {
				newProfile[i] = append(newProfile[i][:j-numberOfDeletions], newProfile[i][j+1-numberOfDeletions:]...)
				numberOfDeletions++
			}
		}
	}
	return
}

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	alts := make([]Alternative, len(p[0]))
	copy(alts, p[0])

	nbVictoires := make(Count)
	tiebreakFunc := TieBreakFactory(alts)
	scfFunc := SCFFactory(MajoritySCF, tiebreakFunc)

	// On itère sur les alternatives (sauf la dernière car elle aura déjà fait tous les duels)
	for i := 0; i < len(alts)-1; i++ {
		// On itère sur les alternatives qui n'ont pas encore été comparées à i
		for j := i + 1; j < len(alts); j++ {
			bestAlt, err := scfFunc(keepOnlyTwoAlternatives(p, alts[i], alts[j]))
			if err != nil {
				return alts, err
			}
			nbVictoires[bestAlt] = nbVictoires[bestAlt] + 1
		}
	}
	for alt, nbVictoire := range nbVictoires {
		if nbVictoire == len(alts)-1 {
			bestAlts = append(bestAlts, alt)
			return
		}
	}

	return
}
