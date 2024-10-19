package selector

func (s selection) getOtherRankFunction() rankFunc {
	return s.rankBySeeders
}
