package util

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

func DetectGPU() (string, error) {
	zap.L().Debug("Trying to detect gpu for hardware acceleration")

	files, err := os.ReadDir("/dev/dri")
	if err != nil {
		return "", fmt.Errorf("failed to read files in /dev/dri, %w", err)
	}

	for _, f := range files {
		if !strings.HasPrefix(f.Name(), "card") {
			continue
		}

		vendorPath := fmt.Sprintf("/sys/class/drm/%s/device/vendor", f.Name())

		data, err := os.ReadFile(vendorPath)
		if err != nil {
			return "", fmt.Errorf("failed to read vendor file, %w", err)
		}

		v := strings.TrimSpace(string(data))
		switch v {
		case "0x10de":
			return "nvidia", nil
		case "0x8086":
			return "intel", nil
		case "0x1002":
			return "amd", nil
		}
	}

	return "", nil
}
