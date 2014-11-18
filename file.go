package qtrim

import (
	"bufio"
	"hyraxbio.kilnhg.com/golang/bioutil.git"
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
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	input := bufio.NewReader(inputFile)
	defer func() {
		inputFile.Close()
	}()
	return TrimIO(input, output, mean, window, minLength, trimN)
}

func TrimIO(input *bufio.Reader, output *bufio.Writer, mean int, window int, minLength int, trimN bool) ([]interface{}, error) {
	results, err := bioutil.ScanFastq(input, func(read *bioutil.Read) (interface{}, error) {
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

func TrimFileFilter(inputPath string, outputPath string, filter bioutil.SequenceFilterFunc, mean int, window int, minLength int, trimN bool) ([]interface{}, error) {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	output := bufio.NewWriter(outputFile)
	defer func() {
		output.Flush()
		outputFile.Close()
	}()
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	input := bufio.NewReader(inputFile)
	defer func() {
		inputFile.Close()
	}()
	return TrimIOFilter(input, output, filter, mean, window, minLength, trimN)
}

func TrimIOFilter(input *bufio.Reader, output *bufio.Writer, filter bioutil.SequenceFilterFunc, mean int, window int, minLength int, trimN bool) ([]interface{}, error) {
	results, err := bioutil.ScanFastqFilter(input, filter, func(read *bioutil.Read) (interface{}, error) {
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
