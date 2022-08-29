package utils

import "gitlab.com/snap-clickstaff/go-common/comutils"

func BytesSlice(bytes []byte, from, to int) []byte {
	byteSize := len(bytes)

	if from < 0 {
		from = (byteSize + from) % byteSize
	} else {
		from = comutils.MinInt(from, byteSize)
	}
	if to < 0 {
		to = (byteSize + to) % byteSize
	} else {
		to = comutils.MinInt(to, byteSize)
	}

	return bytes[from:to]
}
