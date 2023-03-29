package shuttlefolder

import (
	"encoding/hex"
	"fmt"
	"path"
)

const (
	TaskBinaryDir    string = "binaries"
	TaskBinaryPrefix        = "shuttletask"
)

func CalculateBinaryPath(shuttledir, hash string) string {
	return path.Join(
		shuttledir,
		"binaries",
		fmt.Sprintf("%s-%s", TaskBinaryPrefix, hex.EncodeToString([]byte(hash)[:16])),
	)
}
