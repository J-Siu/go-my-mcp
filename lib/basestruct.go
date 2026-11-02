package lib

import "errors"

// Embedded - common fields to be embedded by other structs
type Embedded struct {
	Debug       bool   `json:"Debug"` // simple log level control
	Err         error  `json:"Err,omitempty"`
	Initialized bool   `json:"Initialized,omitempty"`
	LogLevel    uint8  `json:"LogLevel,omitempty"` // finer log control
	MyType      string `json:"MyType"`             // Store typename. Cheaper way than reflector for logging.
}

// IsInitAndNoErr - check error and initialization state in following order:
//  1. If `Err` not nil -> return `false`
//  2. If `Initialized` is `false` -> set `Err` -> return `false`
//  3. Else, check passed -> return `true`
func (t *Embedded) IsInitAndNoErr(prefix string) (pass bool) {
	pass = true
	if t.Err != nil {
		pass = false
	} else if !t.Initialized {
		errMsg := "not initialized"
		if prefix != "" {
			errMsg = prefix + ": " + errMsg
		} else if t.MyType != "" {
			errMsg = t.MyType + ": " + errMsg
		}
		t.Err = errors.New(errMsg)
		pass = false
	}
	return pass
}

// GetErr - for interface
func (t *Embedded) GetErr() error { return t.Err }

// ErrNew - create an error prefixed by `<prefix>: ` when prefix is not empty
func ErrNew(prefix string, strErr string) error {
	if len(prefix) == 0 {
		return errors.New(strErr)
	}
	return errors.New(prefix + ": " + strErr)
}
