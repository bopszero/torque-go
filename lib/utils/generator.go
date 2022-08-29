package utils

import (
	"io"
)

type (
	ReaderGeneratorLoader func() ([]byte, error)
	ReaderGenerator       struct {
		loader      ReaderGeneratorLoader
		cachedBytes []byte
		isExhausted bool
	}
)

func NewReaderGenerator(loader ReaderGeneratorLoader) *ReaderGenerator {
	return &ReaderGenerator{
		loader: loader,
	}
}

func (this *ReaderGenerator) Read(target []byte) (n int, err error) {
	if len(target) == 0 {
		return 0, IssueErrorf("ReaderGenerator.Read cannot read to empty target")
	}
	if this.isExhausted {
		return 0, io.EOF
	}
	for !this.isExhausted && len(this.cachedBytes) < len(target) {
		newBytes, err := this.loader()
		if err != nil {
			return 0, err
		}
		if len(newBytes) == 0 {
			this.isExhausted = true
			break
		}
		this.cachedBytes = append(this.cachedBytes, newBytes...)
	}
	if len(this.cachedBytes) == 0 {
		return 0, io.EOF
	}
	count := copy(target, this.cachedBytes)
	this.cachedBytes = this.cachedBytes[count:]
	return count, nil
}
