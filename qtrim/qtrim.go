package main

import (
	"bufio"
	"flag"
	"fmt"
	"code.hyraxbio.co.za/qtrim"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("For usage please type qtrim trim -h\n")
		os.Exit(0)
	}
	switch os.Args[1] {
	case "trim":
		Trim()
	case "plot":
		Plot()
	default:
		panic("Unknown command")
	}
}

func Trim() {
	var inputFile, outputFile string
	var mean, window, minLength int
	var trimN bool
	flags := flag.NewFlagSet("trim", flag.ExitOnError)
	flags.StringVar(&inputFile, "input", "", "Name of the fastq input file. Otherwise STDIN.")
	flags.StringVar(&outputFile, "output", "", "Name of the output file. Otherwise STDOUT.")
	flags.IntVar(&mean, "mean", 20, "Mean quality of output sequences.")
	flags.IntVar(&window, "window", 50, "Size of trailing window.")
	flags.IntVar(&minLength, "length", 50, "Minimum sequence length.")
	flags.BoolVar(&trimN, "trimn", true, "Remove trailing N's.")
	flags.Parse(os.Args[2:])
	input := bufio.NewReader(os.Stdin)
	output := bufio.NewWriter(os.Stdout)
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			panic("Cannot open specified input file: " + inputFile)
		}
		input = bufio.NewReader(file)
	}
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			panic("Cannot create specified output file: " + outputFile)
		}
		output = bufio.NewWriter(file)
	}
	qtrim.TrimIO(input, output, mean, window, minLength, trimN)
}

func Plot() {
	panic("not implemented yet!!")
}
