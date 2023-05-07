package helpers

import (
	"fmt"
	"strings"
)

func (m *MultiError) Error() string {
	var errorMessages []string
	for _, err := range m.Errors {
		errorMessages = append(errorMessages, err.Error())
	}
	return strings.Join(errorMessages, "; ")
}

func buildInputs(manifestData Manifest, userInput map[string]interface{}) (map[string]interface{}, error) {
	inputs := make(map[string]interface{})
	var multiError MultiError

	for _, input := range manifestData.Inputs {
		// Check if the input value is provided by the user
		if value, ok := userInput[input.Name]; ok {
			// If the input is required, and the user provided a value, use it
			inputs[input.Name] = value
		} else {
			// If the input is required and not provided by the user, add an error
			if input.Required {
				multiError.Errors = append(multiError.Errors, fmt.Errorf("input %s is required, but no value provided", input.Name))
			} else {
				// If the input is not required, use the default value if provided
				inputs[input.Name] = input.Default
			}
		}
	}

	if len(multiError.Errors) > 0 {
		return nil, &multiError
	}
	return inputs, nil
}
