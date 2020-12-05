package fuzzers

// Fuzz input must be of form `pid Read` or `pid Write val`
// A return value of 1 denotes interesting input and thus non-linearizable histories should produce 1.
func Fuzz(data []byte) int {
	// if _, err := Decode(bytes.NewReader(data)); err != nil {
	s := string(data)
	isValid := isValidInput(s)
	if isValid == -1 {
		// input cannot be parsed.
		return -1
	}
	// TODO(rohailasim123): Write code to run tests with input `s`. See corpus
	// for format of input. Return 0 if history is linearizable, 1 otherwise.

}
