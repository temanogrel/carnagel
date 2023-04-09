package encoder

import (
	"github.com/pkg/errors"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func ReplaceFileExtension(source, ext string) string {
	oldExt := filepath.Ext(source)

	if strings.HasPrefix(ext, ".") {
		return source[0:len(source)-len(oldExt)] + ext
	}

	return source[0:len(source)-len(oldExt)] + "." + ext
}

func ExtendFileName(source, suffix, ext string) string {

	filename := filepath.Base(source)
	n := strings.LastIndexByte(filename, '.')
	if n >= 0 {
		filename = filename[:n]
	}

	return filepath.Join(
		filepath.Dir(source),
		filename+suffix+"."+ext,
	)
}

func DurationToInterval(duration float64) uint {
	var generateImages float64
	switch {
	case duration < 10*60:
		generateImages = 16

	case duration < 30*60:
		generateImages = 32

	default:
		generateImages = 64
	}

	return uint(math.Ceil(duration / (generateImages - 1)))
}

func HasEnoughSpace(minSpace uint64) (bool, error) {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(os.TempDir(), &stat); err != nil {
		return false, errors.Wrapf(err, "Failed to retrieve available space")
	}

	// Available blocks * size per block = available space in bytes
	return (stat.Bavail * uint64(stat.Bsize)) >= minSpace, nil
}

// FilesSortedByNumber sorts the file according to a specific pattern
// prefix-:number.jpg
type FilesSortedByNumber []string

func (a FilesSortedByNumber) Len() int      { return len(a) }
func (a FilesSortedByNumber) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FilesSortedByNumber) Less(i, j int) bool {

	file1 := filepath.Base(a[i])
	file2 := filepath.Base(a[j])

	getNumber := func(filename string) int {
		parts := strings.Split(filename, "-")

		if len(parts) != 2 {
			// Yes, felix I know
			panic("Invalid pattern expected prefix-<number>.ext")
		}

		parts = strings.Split(parts[1], ".")

		if len(parts) != 2 {
			// Yes, felix I know
			panic("Invalid pattern expected prefix-<number>.ext")
		}

		if number, err := strconv.Atoi(parts[0]); err == nil {
			return number
		} else {
			panic(err)
		}
	}

	return getNumber(file1) < getNumber(file2)
}
