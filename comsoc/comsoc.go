package comsoc

import (
	"errors"
	"fmt"
)

type Alternative int
type Profile [][]Alternative
type Count map[Alternative]int

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int {
	for i, prefAlt := range prefs {
		if prefAlt == alt {
			return i
		}
	}
	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	a := rank(alt1, prefs)
	b := rank(alt2, prefs)
	return b == -1 || (a != -1 && a < b)
}

// renvoie les meilleures alternatives pour un décomtpe donné
func maxCount(count Count) (bestAlts []Alternative) {
	max := 0
	// On recherche le nombre de votes maximal
	for alternative, value := range count {
		if value > max {
			max = value
			bestAlts = make([]Alternative, 0, 1)
		}
		if value == max {
			bestAlts = append(bestAlts, alternative)
		}
	}

	return
}

// vérifie les préférences d'un agent, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois
func checkProfile(prefs []Alternative, alts []Alternative) error {
	if len(prefs) < len(alts) {
		return errors.New("an alternative missing in profile")
	}

	for _, alt := range alts {
		found := false
		for _, alt2 := range prefs {
			if alt == alt2 {
				found = true
				break
			}
		}
		if !found {
			return errors.New("an alternative is present multiple times in profile")
		}
	}

	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error {
	for i, profile := range prefs {
		if err := checkProfile(profile, alts); err != nil {
			return fmt.Errorf("profile %v : %s", i, err.Error())
		}
	}
	return nil
}

func checkAlternativesUnicity(alts []Alternative) error {

	altCounts := make(map[Alternative]int)

	for _, alt := range alts {
		if altCounts[alt] > 1 {
			return fmt.Errorf("the alternative %v is present twice in alternatives", alt)
		}
		altCounts[alt] = 1
	}

	return nil
}
