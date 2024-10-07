package comsoc

func BordaSWF(p Profile) (count Count, err error) {
	var alts []Alternative
	copy(alts, p[0])

	if err = checkAlternativesUnicity(alts); err != nil {
		return
	}
	if err = checkProfileAlternative(p, alts); err != nil {
		return
	}

	count = make(Count)

	var n int
	for _, voter := range p {
		n = len(voter) - 1
		for _, alt := range voter {
			count[alt] = count[alt] + n
			n = n - 1
		}
	}

	return
}

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)
	if err != nil {
		return
	}
	bestAlts = maxCount(count)
	return
}
