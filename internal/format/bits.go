package format

import (
	"fmt"
)

var units = []string{"Bits", "Kilobits", "Megabits", "Gigabits", "Terabits",
	"Petabits", "Exabits", "Zettabits", "Yottabits", "Brontobytes"}

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

func FormatBits(b float64) string {
	idx := 0
	for b >= 1000 && idx < len(units)-1 {
		b /= 1000
		idx++
	}
	return fmt.Sprintf("%.2f %s", b, units[idx])
}

func FormatBytes(bits float64) string {
	b := bits / 8
	idx := 0
	for b >= 1024 && idx < len(byteUnits)-1 {
		b /= 1024
		idx++
	}
	return fmt.Sprintf("%.2f %s", b, byteUnits[idx])
}
