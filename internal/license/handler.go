package license

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bootengine/boot/internal/assets"
)

var availableLicenses map[string]string = map[string]string{
	"mit":           "MIT.License",
	"gnugpl3":       "GPL.LICENSE",
	"gnuagpl3":      "AGPL.LICENSE",
	"gnulgpl3":      "LGPL.LICENSE",
	"mozillapublic": "MOZILLA.LICENSE",
	"apache2":       "APACHE.LICENSE",
	"boostsoftware": "BOOT.LICENSE",
	"unlicense":     "UNLICENSE.LICENSE",
}

func GetLicenseContent(ctx context.Context, licenseName string) (*string, error) {
	filename, ok := availableLicenses[licenseName]
	if !ok {
		return nil, fmt.Errorf("license %s is not available. please create an issue or consider contributing.", licenseName)
	}

	licenseFile, err := assets.LicenseFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content := string(licenseFile)
	// template in some case
	switch licenseName {
	case "mit":
		year := time.Now().Year()
		name := ctx.Value("") // TODO: get value from context or error out

		content = strings.ReplaceAll(content, "[year]", strconv.Itoa(year))
		content = strings.ReplaceAll(content, "[fullname]", name.(string))
	case "apache":
		year := time.Now().Year()
		name := ctx.Value("") // TODO: get value from context or error out

		content = strings.ReplaceAll(content, "[yyyy]", strconv.Itoa(year))
		content = strings.ReplaceAll(content, "[name of copyright owner]", name.(string))
	}

	return &content, nil
}
