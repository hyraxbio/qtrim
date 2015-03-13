package qtrim

import (
	"code.hyraxbio.co.za/bioutil"
)

var base byte = []byte("!")[0]

// Trim a single line.
// Returns the length of the valid sequence, or 0 to discard.
// trimmed = original[0:cutoff - 1]
func Trim(read bioutil.Read, mean int, window int, minLength int) (cutoff int, originalMean float32, finalMean float32) {
	if window > minLength {
		panic("QTrim: window must be <= minimum read length")
	}

	line := read.Quality()
	length := read.Length()

	total := 0
	windowTotal := 0
	for i, b := range line {
		quality := Score(b)
		total += quality
		if i > length-window {
			windowTotal += quality
		}
	}

	originalMean = float32(total) / float32(length)

	windowCompare := mean * window
	cutoff = length
	if cutoff <= minLength || cutoff <= window { // Discard this sequence
		return 0, originalMean, 0
	}
	for !(total >= mean*cutoff && windowTotal >= windowCompare) {
		if cutoff <= minLength || cutoff <= window { // Discard this sequence
			return 0, originalMean, 0
		}
		cutoff--
		trimmedScore := Score(line[cutoff])
		total -= trimmedScore
		windowTotal = windowTotal - trimmedScore + Score(line[cutoff-window])
	}
	for Score(line[cutoff-1]) < mean {
		if cutoff <= minLength || cutoff <= window { // Discard this sequence
			return 0, originalMean, 0
		}
		cutoff--
	}
	finalMean = float32(total) / float32(cutoff)
	return
}

func Score(letter byte) int {
	return int(letter) - int(base)
}
