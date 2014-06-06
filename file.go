package qtrim

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Stat struct {
	InputLength  int
	OutputLength int
	OriginalMean float32
	FinalMean    float32
}

type Read struct {
	HeadLine []byte
	SeqLine  []byte
	SepLine  []byte
	QualLine []byte
}

type SequenceFunc func(read Read) (interface{}, error)

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
	results, err := scanFile(inputPath, func(read Read) (interface{}, error) {
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

func scanFile(inputPath string, function SequenceFunc) ([]interface{}, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	input := bufio.NewReader(inputFile)
	defer func() {
		inputFile.Close()
	}()
	info, _ := inputFile.Stat()
	results := make([]interface{}, 0, info.Size()/400)
	for headLine, err := input.ReadBytes('\n'); err != io.EOF; headLine, err = input.ReadBytes('\n') {
		seqLine, _ := input.ReadBytes('\n')
		sepLine, _ := input.ReadBytes('\n')
		qualLine, _ := input.ReadBytes('\n')
		result, err := function(Read{headLine, seqLine, sepLine, qualLine})
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
