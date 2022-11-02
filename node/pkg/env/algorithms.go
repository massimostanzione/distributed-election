// Algorithms involved in the application and utility functions.
package env

type DEAlgorithm uint8

const (
	DE_ALGORITHM_UNDEFINED DEAlgorithm = iota
	DE_ALGORITHM_BULLY
	DE_ALGORITHM_FREDRICKSONLYNCH
)

func ParseDEAlgorithm(input string) DEAlgorithm {
	bullyFlags := []string{"BULLY", "b", "DE_ALGORITHM_BULLY"}
	flFlags := []string{"FREDRICKSONLYNCH", "fl", "DE_ALGORITHM_FREDRICKSONLYNCH"}
	if contains(bullyFlags, input) {
		return DE_ALGORITHM_BULLY
	}
	if contains(flFlags, input) {
		return DE_ALGORITHM_FREDRICKSONLYNCH
	}
	return DE_ALGORITHM_UNDEFINED
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
