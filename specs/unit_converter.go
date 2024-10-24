package specs

// Assumes mega to be 10^6 as defined by S.I. Function rounds down.
func byteToMegabyte(bytesAmount uint64) uint64 {
	return bytesAmount / 1000 / 1000
}

func kiloHertzToMegaHertz(bytesAmount uint64) uint64 {
	return bytesAmount / 1000
}
