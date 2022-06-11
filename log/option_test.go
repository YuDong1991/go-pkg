package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Options_Validate(t *testing.T) {
	option := &Options{
		Level:            "test",
		Format:           "test",
		EnableColor:      true,
		DisableCaller:    false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	errs := option.Validate()
	expected := `[unrecognized level: "test" not a valid log format: "test"]`
	fmt.Printf("%s\n", errs)
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))
}
