package license

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bootengine/boot/internal/assets"
	"github.com/bootengine/boot/internal/helper"
)

var availableLicenses map[string]string = map[string]string{
	"mit":           "MIT.LICENSE",
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
		return nil, fmt.Errorf("license %s is not available. please create an issue or consider contributing", licenseName)
	}

	licenseFile, err := assets.LicenseFS.ReadFile(filepath.Join("licenses", filename))
	if err != nil {
		return nil, err
	}

	content := string(licenseFile)
	if licenseName == "mit" || licenseName == "apache" {
		year := time.Now().Year()
		name := ctx.Value(helper.ValueKey{}).(map[string]any)["owner"]
		if name == nil {
			return nil, fmt.Errorf("for the selected license (%s), an 'owner' must be defined using vars in the config file", licenseName)
		}
		var yearString, nameString string

		if licenseName == "mit" {
			yearString = "[year]"
			nameString = "[fullname]"
		}
		if licenseName == "apache" {
			yearString = "[yyyy]"
			nameString = "[name of copyright owner]"
		}
		content = strings.ReplaceAll(content, yearString, strconv.Itoa(year))
		content = strings.ReplaceAll(content, nameString, name.(string))
	}

	return &content, nil
}
