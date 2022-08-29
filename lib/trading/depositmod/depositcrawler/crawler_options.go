package depositcrawler

import (
	"gitlab.com/snap-clickstaff/torque-go/config"
)

type CrawlerOptions struct {
	BlockScanBackward uint64
	BlockScanForward  uint64
}

func (this CrawlerOptions) GenBlockHeights(indexHeight uint64) (heights []uint64) {
	var (
		lastBlockHeight   = indexHeight + this.BlockScanForward
		blockScanBackward = this.BlockScanBackward
	)
	if config.Test {
		blockScanBackward /= 10
	}

	heights = make([]uint64, 0, int(blockScanBackward+1+this.BlockScanForward))
	for height := indexHeight - blockScanBackward; height <= lastBlockHeight; height++ {
		heights = append(heights, height)
	}
	return heights
}
