package result

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func NewError(err string, code string) *Response {
	return &Response{
		ExitStatus: 1,
		Error:      err,
		ErrorCode:  code,
	}
}

func NewErrorWithDetails(err string, code string, details error) *Response {
	detailsString := ""

	if details != nil {
		detailsString = details.Error()
	}

	return &Response{
		ExitStatus: 1,
		Error:      err,
		ErrorCode:  code,
		Details:    detailsString,
	}
}

func NewSuccess(success string) *Response {
	return &Response{
		ExitStatus: 0,
		Success:    success,
	}
}

func NewSerializable(success string, serializable any) *Response {
	return &Response{
		ExitStatus:   0,
		Success:      success,
		Serializable: serializable,
	}
}

// Returned by command functions
type Response struct {
	Error     string `json:"error,omitempty"`         // hand written error message
	Details   string `json:"error_details,omitempty"` // underlying error message
	ErrorCode string `json:"code,omitempty"`          // machine-keyable error code
	Success   string `json:"success,omitempty"`       // hand written success message

	Serializable any   `json:"-"` // wholesale output if mode is JSON
	ExitStatus   uint8 `json:"-"` // process exit status code. keeping hidden since it's not always present

	SuccessMultiline bool `json:"-"` // success message is multi-line, disable automatic colors
	ErrorMultiline   bool `json:"-"` // error message is multi-line, disable automatic colors
}

func (r *Response) AddSuccessLn(line string) {
	// when switching from single to multi line, assuming we already have content, add a newline
	if !r.SuccessMultiline && r.Success != "" {
		r.Success += "\n"
	}

	r.Success += line + "\n"

	r.SuccessMultiline = true
}

func (r *Response) AddErrorLn(line string) {
	// when switching from single to multi line, assuming we already have content, add a newline
	if !r.ErrorMultiline && r.Error != "" {
		r.Error += "\n"
	}

	r.Error += line + "\n"

	r.ErrorMultiline = true
}

func (r *Response) SetError(err string, code string) {
	r.ExitStatus = 1
	r.Error = err
	r.ErrorCode = code
}

func (r *Response) SetErrorDetails(details error) {
	r.ExitStatus = 1
	r.Details = details.Error()
}

func (r Response) Display(encode bool) error {
	if encode {
		if r.Serializable != nil {
			serialized, err := json.Marshal(r.Serializable)
			if err != nil {
				fmt.Println("unable to serialize custom output json")
				return err
			}

			fmt.Println(string(serialized))
		} else {
			serialized, err := json.Marshal(r)
			if err != nil {
				fmt.Println("unable to serialize response output json")
				return err
			}

			fmt.Println(string(serialized))
		}
	} else {
		if r.Error != "" {
			if r.ErrorMultiline {
				fmt.Println(strings.TrimSuffix(r.Error, "\n"))
			} else {
				color.Red(strings.TrimSuffix(r.Error, "\n"))
			}
		}

		if r.Details != "" {
			color.Yellow(strings.TrimSuffix(r.Details, "\n"))
		}

		if r.Success != "" {

			if r.ErrorMultiline {
				fmt.Println(strings.TrimSuffix(r.Success, "\n"))
			} else {
				color.HiGreen(strings.TrimSuffix(r.Success, "\n"))
			}
		}
	}

	return nil
}
