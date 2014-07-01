package qtrim

import (
	"fmt"
	"github.com/SANBIHIV/biotools/bioutil"
	"github.com/baruchlubinsky/plotly"
)

const bins = 40

const QUALITYLAYOUT = `{
  "title":"Distribution of quality scores",
  "yaxis":{
    "title": "Quality score",
    "zeroline": true,
    "overlaying":false,
    "side":"left",
    "anchor":"x"
  },
  "yaxis2":{
    "title":"Coverage",
    "overlaying":"y",
    "side":"right",
    "anchor":"x"
  },
  "xaxis":{
    "title":"Read position"
  },
  "showlegend":false
}`

func BoxAndWhisker(inputPath string, binSize int) plotly.Figure {
	boxes := make([][]int, 0, 400/binSize)
	totals := make([]int, 0, cap(boxes))
	bioutil.ScanFastqFile(inputPath, func(read bioutil.Read) (interface{}, error) {
		for i := 0; i < len(read.QualLine)-1; i = i + binSize {
			bin := i / binSize
			for len(boxes) < bin+1 {
				boxes = append(boxes, make([]int, bins+1))
				totals = append(totals, 0)
			}
			score := int(Score(read.QualLine[i]))
			if score < 0 {
				continue
			}
			if score > bins {
				score = bins
			}
			boxes[bin][int(score)]++
			totals[bin]++
		}
		return nil, nil
	})
	result := make([]plotly.Trace, len(boxes))
	for i, _ := range result {
		points := 20.0
		result[i].Y = make([]interface{}, int(points))
		result[i].X = make([]interface{}, int(points))
		result[i].Name = fmt.Sprintf("%v", i*binSize)
		result[i].YAxis = "y"
		result[i].Type = "box"
		result[i].Line = plotly.Line{
			Color: "rgb(0.1, 0.3, 0.9)",
		}
		result[i].Marker = plotly.Marker{
			Line:    plotly.Line{},
			Opacity: 0.0,
		}
		data := boxes[i]
		accumulation := 0
		j := 0
		total := float64(totals[i])
		for score := 0; score <= 40; score++ {
			accumulation += data[score]
			for accumulation > 0 && float64(accumulation)/total > (1.0/points)*float64(j) {
				result[i].Y[j] = float64(score)
				result[i].X[j] = float64(i * binSize)
				j++
				if float64(j) > points-1 {
					break
				}
			}
			if float64(j) > points-1 {
				break
			}
		}
	}
	counts := plotly.Trace{
		Y:     make([]interface{}, len(totals)),
		X:     make([]interface{}, len(totals)),
		Name:  "Coverage",
		YAxis: "y2",
		Type:  "scatter",
		Line: plotly.Line{
			Color: "rgb(0.2, 0.8, 0.1)",
		},
		Marker: plotly.Marker{
			Line:    plotly.Line{},
			Opacity: 1.0,
		},
	}
	for i, count := range totals {
		counts.X[i] = i * binSize
		counts.Y[i] = count
	}
	result = append(result, counts)
	return plotly.Figure{
		Data:   result,
		Layout: QUALITYLAYOUT,
	}
}

func QualityTrendPlot(input string, name string, output string) error {
	fig := BoxAndWhisker(input, 10)
	url, err := plotly.Create(name, fig, true)
	if err != nil {
		return err
	}
	err = plotly.Save(url.Id(), output)
	return err
}
