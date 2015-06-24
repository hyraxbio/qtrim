package qtrim

import (
	"bufio"
	"code.hyraxbio.co.za/bioutil"
	"code.hyraxbio.co.za/bioutil/dnaio"
	"os"
)

type Stat struct {
	InputLength  int
	OutputLength int
	OriginalMean float32
	FinalMean    float32
}

func TrimFile(inputPath string, outputPath string, mean int, window int, minLength int) ([]Stat, error) {
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
	return TrimIO(input, output, mean, window, minLength), nil
}

func TrimIO(input *bufio.Reader, output *bufio.Writer, mean int, window int, minLength int) []Stat {
	inputChan := make(chan bioutil.Read)
	outputChan := make(chan bioutil.Read)
	go bioutil.ScanFastqChan(input, inputChan)
	go func() {
		for read := range outputChan {
			output.Write(read.Data())
		}
	}()

	results := TrimPipe(inputChan, outputChan, mean, window, minLength)

	return results
}

func TrimPipe(input chan bioutil.Read, output chan bioutil.Read, mean int, window int, minLength int) []Stat {
	results := make([]Stat, 0)
	for read := range input {

		length := len(read.Sequence())
		count, original, final := Trim(read, mean, window, minLength)
		result := Stat{
			length,
			count,
			original,
			final,
		}
		if count != 0 {
			trimmed := length - count
			trimmedRead := read.TrimRight(trimmed)
			output <- trimmedRead
		}
		results = append(results, result)
	}
	close(output)
	return results
}
