// version 2.0.0

package comsoc

import "testing"

func TestBordaSWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := BordaSWF(prefs)

	if res[1] != 4 {
		t.Errorf("error, result for 1 should be 4, %d computed", res[1])
	}
	if res[2] != 3 {
		t.Errorf("error, result for 2 should be 3, %d computed", res[2])
	}
	if res[3] != 2 {
		t.Errorf("error, result for 3 should be 2, %d computed", res[3])
	}
}

func TestBordaSCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := BordaSCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestMajoritySWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := MajoritySWF(prefs)

	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 0 {
		t.Errorf("error, result for 2 should be 0, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestMajoritySCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, err := MajoritySCF(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestApprovalSWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 3, 2},
		{2, 3, 1},
	}
	thresholds := []int{2, 1, 2}

	res, _ := ApprovalSWF(prefs, thresholds)
	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 2 {
		t.Errorf("error, result for 2 should be 2, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestApprovalSCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
	}
	thresholds := []int{2, 1, 2}

	res, err := ApprovalSCF(prefs, thresholds)

	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestTieBreakFactory(t *testing.T) {
	prefs := [][]Alternative{
		{1, 3, 2},
		{2, 1, 3},
		{2, 1, 3},
	}
	thresholds := []int{2, 1, 2}

	res, err := ApprovalSCF(prefs, thresholds)
	if err != nil {
		t.Error(err)
	}
	tieBreak := TieBreakFactory(prefs[0])
	result, err := tieBreak(res)
	if err != nil {
		t.Error(err)
	}
	if result != 1 {
		t.Errorf("Erreur, le gagnat devrait être 1")
	}
}

func TestSCFFactory(t *testing.T) {
	prefs := [][]Alternative{
		{1, 3, 2},
		{2, 1, 3},
		{2, 1, 3},
	}
	thresholds := []int{2, 1, 2}
	appro2 := func(p Profile) (alt []Alternative, err error) {
		return ApprovalSCF(p, thresholds)
	}
	tieBreak := TieBreakFactory(prefs[0])
	fn := SCFFactory(appro2, tieBreak)
	result, err := fn(prefs)
	if err != nil {
		t.Error(err)
	}
	if result != 1 {
		t.Errorf("Erreur, le gagnat devrait être 1")
	}
}

func TestSWFFactory(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{2, 1, 3},
		{2, 1, 3},
		{2, 3, 1},
		{1, 3, 1},
		{3, 2, 1},
	}

	trieur := []Alternative{1, 3, 2}
	tieBreak := TieBreakFactory(trieur)
	fn := SWFFactory(MajoritySWF, tieBreak)
	//t.Error(MajoritySWF(prefs))
	res, err := fn(prefs)

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 {
		t.Errorf("list should be of size 1, not %v", len(res))
	}

	if res[0] != 1 {
		t.Errorf("the winner should be 1, not %v", res[0])
	}
}

func TestCondorcetWinner(t *testing.T) {
	prefs1 := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	prefs2 := [][]Alternative{
		{1, 2, 3},
		{2, 3, 1},
		{3, 1, 2},
	}

	res1, _ := CondorcetWinner(prefs1)
	res2, _ := CondorcetWinner(prefs2)

	if len(res1) == 0 || res1[0] != 1 {
		t.Errorf("error, 1 should be the only best alternative for prefs1. Got %v", res1)
	}
	if len(res2) != 0 {
		t.Errorf("no best alternative for prefs2")
	}
}
