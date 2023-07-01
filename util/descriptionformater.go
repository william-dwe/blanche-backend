package util

import "fmt"

const (
	TOPUP_DESCRPTION = "Top Up from"
)

func TopUpDescFormater(sourceFund string) string {
	return fmt.Sprintf("%s %s", TOPUP_DESCRPTION, sourceFund)
}
