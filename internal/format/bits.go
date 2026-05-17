package format

import (
	"fmt"
)

var units = []string{"Bits", "Kilobits", "Megabits", "Gigabits", "Terabits",
	"Petabits", "Exabits", "Zettabits", "Yottabits", "Brontobytes"}

func FormatBits(b float64) string {
	idx := 0
	for b >= 1000 && idx < len(units)-1 {
		b /= 1000
		idx++
	}
	return fmt.Sprintf("%.2f %s", b, units[idx])
}
