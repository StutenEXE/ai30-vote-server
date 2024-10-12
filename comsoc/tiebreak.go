package comsoc

import (
	"errors"
	"sort"
)

func TieBreakFactory(orderedAlts []Alternative) func([]Alternative) (Alternative, error) {
	return func(alts []Alternative) (winner Alternative, err error) {
		for _, alt := range orderedAlts {
			for _, possibleWinner := range alts {
				if possibleWinner == alt {
					return possibleWinner, nil
				}
			}
		}
		return alts[0], errors.New("erreur, pas de gagant dans la liste")
	}
}

func SWFFactory(swf func(p Profile) (Count, error), tb func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	return func(p Profile) (alts []Alternative, err error) {
		c, err := swf(p)
		if err != nil {
			return
		}
		for elem := range c {
			alts = append(alts, elem)
		}
		alts = maxCount(c)
		sort.Slice(alts, func(i, j int) bool {
			winner, _ := tb([]Alternative{alts[i], alts[j]})
			return alts[i] == winner
		})
		return
	}
}

func SCFFactory(scf func(p Profile) ([]Alternative, error), tb func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	return func(p Profile) (alt Alternative, err error) {
		c, err := scf(p)
		if err != nil {
			return
		}
		alt, err = tb(c)
		return
	}
}
