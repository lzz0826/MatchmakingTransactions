package utils

import "fmt"

func SafeDiv(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("除數不能為 0")
	}
	return a / b, nil
}
