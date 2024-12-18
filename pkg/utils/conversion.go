package utils

// given days return it in seconds
func ConvertDaysInSeconds(days uint64) uint64 {
	return 60 * 60 * 24 * days
}
