package comsoc

import "errors"

func MajoritySWF(p Profile) (count Count, err error) {
	if len(p) == 0 {
		err = errors.New("empty profile")
		return
	}

	var alts []Alternative
	copy(alts, p[0])

	if err = checkAlternativesUnicity(alts); err != nil {
		return
	}
	if err = checkProfileAlternative(p, alts); err != nil {
		return
	}

	count = make(Count)
	for _, votant := range p {
		alt := votant[0]
		count[alt] = count[alt] + 1
	}
	return
}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	var c Count
	c, err = MajoritySWF(p)
	bestAlts = maxCount(c)
	return
}
