package qtrim

import (
	"bufio"
	"github.com/SANBIHIV/biotools/bioutil"
	"os"
	"strings"
)

type Stat struct {
	InputLength  int
	OutputLength int
	OriginalMean float32
	FinalMean    float32
}

func TrimFile(inputPath string, outputPath string, mean int, window int, minLength int, trimN bool) ([]interface{}, error) {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	output := bufio.NewWriter(outputFile)
	defer func() {
		output.Flush()
		outputFile.Close()
	}()
	results, err := bioutil.ScanFastqFile(inputPath, func(read bioutil.Read) (interface{}, error) {
		if trimN {
			read.SeqLine = []byte(strings.TrimRight(string(read.SeqLine), "N\n"))
		}
		length := len(read.SeqLine)
		count, original, final := Trim(read.QualLine[0:length], mean, window, minLength)
		result := Stat{
			length,
			count,
			original,
			final,
		}
		if count != 0 {
			output.Write(read.HeadLine)
			output.Write(read.SeqLine[0:count])
			output.WriteByte('\n')
			output.Write(read.SepLine)
			output.Write(read.QualLine[0:count])
			output.WriteByte('\n')
		}
		return result, nil
	})
	return results, err
}
