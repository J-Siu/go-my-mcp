/*
Copyright © 2026 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package conf

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/J-Siu/go-my-mcp/helper"
	"github.com/spf13/viper"
)

const (
	AppName = "go-my-mcp"

	// Default config file location, used when --FileConf is not given
	DefaultConfigFile = "~/.config/" + AppName + ".json"

	DefaultListenAddress = ":9000"
	DefaultSearXngUrl    = "http://localhost:8080"

	// Maximum number of redirects a single web_fetch follows
	DefaultRedirectMax = 10

	// Response body limits in bytes
	DefaultFetchBodySizeMax  = 2 << 20
	DefaultSearchBodySizeMax = 2 << 20

	// Environment variables are prefixed with this, followed by the upper-cased key name
	EnvPrefix = "GO_MY_MCP"

	// Replaces the API key in the `config` output, so it is never echoed
	txt_Redacted = "<redacted>"
)

// Viper instance backing [Conf]. Cobra flags are bound to it from the `cmd` package.
var Vip = viper.New()

type TypeConf struct {
	ApiKey            string `json:"ApiKey"` // Shared secret required from clients. Empty disables access control.
	Debug             bool   `json:"Debug"`
	FetchBodySizeMax  int64  `json:"FetchBodySizeMax"`
	FileConf          string `json:"FileConf"`
	ListenAddress     string `json:"ListenAddress"`
	RedirectMax       int    `json:"RedirectMax"`
	SearXngUrl        string `json:"SearXngUrl"`
	SearchBodySizeMax int64  `json:"SearchBodySizeMax"`
	SkipTlsVerify     bool   `json:"SkipTlsVerify"`
}

// Runtime configuration, populated by [Conf.Load]
var Conf = TypeConf{
	FetchBodySizeMax:  DefaultFetchBodySizeMax,
	ListenAddress:     DefaultListenAddress,
	RedirectMax:       DefaultRedirectMax,
	SearXngUrl:        DefaultSearXngUrl,
	SearchBodySizeMax: DefaultSearchBodySizeMax,
}

func configKeys() []string {
	confType := reflect.TypeOf(TypeConf{})
	keys := make([]string, 0, confType.NumField())
	for i := 0; i < confType.NumField(); i++ {
		field := confType.Field(i)
		if field.PkgPath == "" && field.Tag.Get("mapstructure") != "-" {
			keys = append(keys, field.Name)
		}
	}
	return keys
}

func init() {
	Vip.SetConfigType("json")
	Vip.SetEnvPrefix(EnvPrefix)
	Vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	conf := reflect.ValueOf(Conf)
	for _, key := range configKeys() {
		Vip.SetDefault(key, conf.FieldByName(key).Interface())
	}
	Vip.AutomaticEnv()

	// keep the logger correct even before the first Load (e.g. init-time errors)
	helper.EnableDebug(false)
}

// Load - resolve the effective configuration. Precedence, lowest to highest:
// config file, environment variables, command line flags.
//
// A missing config file is not an error, the defaults are used instead.
func (t *TypeConf) Load() error {
	const prefix = "TypeConf.Load"
	var (
		err    error
		values TypeConf
	)
	configFile := Vip.GetString("FileConf")
	if configFile == "" {
		configFile = DefaultConfigFile
	}
	Vip.SetConfigFile(helper.TildeEnvExpand(configFile))
	if err = Vip.ReadInConfig(); err != nil {
		if Vip.GetBool("Debug") {
			slog.Warn("could not read config file; continuing without it", "file", configFile, "error", err)
		}
	}
	values = *t
	if err = Vip.Unmarshal(&values); err == nil {
		values.FileConf = configFile
		helper.EnableDebug(values.Debug) // Set the logger before validating
		err = values.validate()
	}
	if err == nil {
		*t = values
	}
	if err != nil {
		err = fmt.Errorf("%s: %w", prefix, err)
	}
	return err
}

// validate - reject values that would break the tools at run time
func (t *TypeConf) validate() error {
	const prefix = "TypeConf.Load"
	switch {
	case t.RedirectMax < 0:
		return fmt.Errorf("%s: RedirectMax must be >= 0, got %d", prefix, t.RedirectMax)
	case t.FetchBodySizeMax <= 0:
		return fmt.Errorf("%s: FetchBodySizeMax must be > 0, got %d", prefix, t.FetchBodySizeMax)
	case t.SearchBodySizeMax <= 0:
		return fmt.Errorf("%s: SearchBodySizeMax must be > 0, got %d", prefix, t.SearchBodySizeMax)
	}
	return nil
}

// SafeConf - create a copy of t with ApiKey redacted if not empty
func SafeConf(t *TypeConf) *TypeConf {
	safe := *t
	if safe.ApiKey != "" {
		safe.ApiKey = txt_Redacted
	}
	return &safe
}

// cspell:words mapstructure
