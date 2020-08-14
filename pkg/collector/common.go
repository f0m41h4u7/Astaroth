package collector

func newMap() map[string]int64 {
	states := map[string]int64{}
	states["LISTEN"] = 0
	states["ESTAB"] = 0
	states["FIN-WAIT"] = 0
	states["SYN-RCV"] = 0
	states["TIME-WAIT"] = 0
	states["CLOSE-WAIT"] = 0

	return states
}
