package xerror

import (
	"bytes"
	"fmt"
)

type InvalidArgsErr struct {
	field   string
	message string
}

func NewInvalidArgsErr(field string, message string) InvalidArgsErr {
	return InvalidArgsErr{
		field:   field,
		message: message,
	}
}

func (e InvalidArgsErr) Error() string {
	return fmt.Sprintf("invalid argument field %q, message: %s", e.field, e.message)
}

type ErrorCollector []error

func (m *ErrorCollector) Error() string {
	buffer := bytes.NewBufferString("")
	for i, err := range *m {
		fmt.Fprintf(buffer, "%d. %s\n", i+1, err.Error())
	}
	return buffer.String()
}

func (m *ErrorCollector) Count() int {
	return len(*m)
}

func (m *ErrorCollector) Collect(e error) {
	*m = append(*m, e)
}
