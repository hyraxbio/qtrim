package qtrim

import (
	"encoding/json"
	"fmt"
)

const bins = 40

type BoxDataPoint struct {
	Name string    `json:"name"`
	Y    []float32 `json:"y"`
}

func BoxAndWhisker(inputPath string, binSize int) string {
	boxes := make([][]int, 0, 400/binSize)
	totals := make([]int, 0, cap(boxes))
	scanFile(inputPath, func(read Read) (interface{}, error) {
		for i, b := range read.QualLine {
			bin := i / 10
			for len(boxes) < bin+1 {
				boxes = append(boxes, make([]int, bins+1))
				totals = append(totals, 0)
			}
			score := Score(b)
			if score < 0 {
				continue
			}
			boxes[bin][int(score)]++
			totals[bin]++
		}
		return nil, nil
	})
	points := []float32{0, 0.25, 0.4, 0.6, 0.75, 0.99}
	result := make([]BoxDataPoint, len(boxes))
	for i, _ := range result {
		result[i].Y = make([]float32, len(points))
		result[i].Name = fmt.Sprintf("%v", i*binSize)
		data := boxes[i]
		accumulation := 0
		j := 0
		for score, count := range data {
			accumulation += count
			if float32(accumulation)/float32(totals[i]) > points[j] {
				result[i].Y[j] = float32(score)
				j++
			}
			if j > len(points)-1 {
				break
			}
		}
	}
	data, _ := json.Marshal(result)
	return string(data)
}
