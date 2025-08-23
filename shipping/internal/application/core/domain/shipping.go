package domain

// Regra: prazo mínimo = 1 dia; a cada 5 unidades, +1 dia.
func EstimateDays(totalUnits int) int32 {
	if totalUnits <= 0 {
		return 1
	}
	extra := (totalUnits - 1) / 5
	return int32(1 + extra)
}
