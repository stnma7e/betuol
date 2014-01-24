package instance

import (
	"github.com/stnma7e/betuol/component"
)

func StartScript(is IInstance, charsToActOn ...component.GOiD) error {
	_, err := is.CreateFromMap("map1")
	if err != nil {
		return err
	}

	return nil
}
