package qtrim

import (
	//"os"
	//"bufio"
	"io"
	"json"
)

const bins = 40

func BoxAndWhisker(inputPath string, binSize int) string {
	boxes := make([][bins]int, 0, 400/binSize)
	totals := make([]int, 0, cap(boxes))
	scanFile(inputPath, func(read Read) (interface{}, error) {
		total := 0
		length := len(read.QualLine)
		for i, b := range read.QualLine {
			bin := i / 10
			score := Score(b)
			total += score
			boxes[bin][int(score)]++
			totals[bin]++
		}
	})
	x := make([]int, len(boxes))
	for i, _ := range boxes {
		x[i] = (i + 1) * binSize
	}
	y := make([][6]int, len(x))

}
