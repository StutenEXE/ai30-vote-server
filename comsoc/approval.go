package comsoc

import "errors"

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	var alts []Alternative
	copy(alts, p[0])

	if err = checkAlternativesUnicity(alts); err != nil {
		return
	}
	if err = checkProfileAlternative(p, alts); err != nil {
		return
	}

	if len(p) != len(thresholds) {
		err = errors.New("error, not as much thersholds as voters")
	}
	count = make(Count)
	for index, nb := range thresholds {
		for i := 0; i < nb; i++ {
			j := p[index][i]
			count[j] = count[j] + 1
		}
	}
	return
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	var c Count
	c, err = ApprovalSWF(p, thresholds)
	bestAlts = maxCount(c)
	return
}
