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

package global

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"fj.vms.local/jsds/go-helper-mini/basestruct"
	"fj.vms.local/jsds/go-helper-mini/errs"
	"github.com/spf13/viper"
)

const (
	AppName = "go-my-mcp"

	// Default config file location, used when --config is not given
	DefaultConfigFile = "~/.config/" + AppName + ".json"

	DefaultListenAddress = ":9000"
	DefaultSearXngUrl    = "http://xng.vms.local"

	// Maximum number of redirects a single web_fetch follows
	DefaultRedirectMax = 10

	// Response body limits in bytes
	DefaultFetchBodySizeMax  = 2 << 20
	DefaultSearchBodySizeMax = 2 << 20

	// Environment variables are prefixed with this, followed by the upper-cased key name
	EnvPrefix = "GO_MY_MCP"
)

const (
	KeyDebug             = "Debug"
	KeyFetchBodySizeMax  = "FetchBodySizeMax"
	KeyFileConf          = "FileConf"
	KeyListenAddress     = "ListenAddress"
	KeyRedirectMax       = "RedirectMax"
	KeySearchBodySizeMax = "SearchBodySizeMax"

	// pseudo key, tracks which flags cobra actually set. Not part of allKeys(),
	// so it is neither reported nor read from the environment.
	KeyConfigFlag = "ConfigFlag"
	KeySearXngUrl = "SearXngUrl"
)

// Effective configuration, populated by [Conf.Load]
var Conf TypeConf

// Viper instance backing [Conf]. Cobra flags are bound to it from the `cmd` package.
var Vip = viper.New()

type TypeConf struct {
	basestruct.Embedded
	FetchBodySizeMax  int64  `json:"FetchBodySizeMax"`
	FileConf          string `json:"FileConf"`
	ListenAddress     string `json:"ListenAddress"`
	RedirectMax       int    `json:"RedirectMax"`
	SearchBodySizeMax int64  `json:"SearchBodySizeMax"`
	SearXngUrl        string `json:"SearXngUrl"`
}

// Values read from config file, environment and flags
type typeConfValues struct {
	ListenAddress string
	SearXngUrl    string
}

// TypeConfEffective - configuration as shown by the `config` sub command
type TypeConfEffective struct {
	FileConf          string            `json:"FileConf"`
	ListenAddress     string            `json:"ListenAddress"`
	SearXngUrl        string            `json:"SearXngUrl"`
	RedirectMax       int               `json:"RedirectMax"`
	FetchBodySizeMax  int64             `json:"FetchBodySizeMax"`
	SearchBodySizeMax int64             `json:"SearchBodySizeMax"`
	Debug             bool              `json:"Debug"`
	Version           string            `json:"Version"`
	Source            map[string]string `json:"Source"`
	Env               TypeConfEnv       `json:"Env"`
}

// Layers a value can come from, in increasing precedence
const (
	sourceDefault = "default"
	sourceEnv     = "env"
	sourceFlag    = "flag"
)

// TypeConfEnv - environment variable state for the configuration keys
type TypeConfEnv struct {
	// Prefix used to build the environment variable names
	Prefix string `json:"Prefix"`
	// Names viper reads, whether set or not
	Recognized []string `json:"Recognized"`
	// Values currently present in the environment
	Set map[string]string `json:"Set"`
}

// Env - environment variables in effect for the configuration
func (t *TypeConf) Env() TypeConfEnv {
	env := TypeConfEnv{
		Prefix: EnvPrefix,
		Set:    map[string]string{},
	}
	for _, key := range allKeys() {
		env.Recognized = append(env.Recognized, envName(key))
	}
	set := map[string]string{}
	for _, entry := range os.Environ() {
		name, value, found := strings.Cut(entry, "=")
		if found && strings.HasPrefix(name, EnvPrefix+"_") {
			set[name] = value
		}
	}
	env.Set = set
	return env
}

// Source - which layer supplied each key: default, env or flag.
//
// Only the winning layer is reported; the config file is folded into `default`
// because a missing file is not an error and not every key is present in it.
// Flag and env are checked before the file, matching [TypeConf.Load] precedence.
func (t *TypeConf) Source() map[string]string {
	source := map[string]string{}
	setFlags, _ := Vip.Get(KeyConfigFlag).(map[string]string)
	for _, key := range allKeys() {
		switch {
		case setFlags[key] != "":
			source[key] = sourceFlag
		case os.Getenv(envName(key)) != "":
			source[key] = sourceEnv
		default:
			source[key] = sourceDefault
		}
	}
	return source
}

// Env var name for a viper key, e.g. `ListenAddress` -> `GO_MY_MCP_LISTENADDRESS`
func envName(key string) string {
	return EnvPrefix + "_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
}

// Viper keys known to the configuration
func allKeys() []string {
	return []string{
		KeyFileConf,
		KeyListenAddress,
		KeySearXngUrl,
		KeyRedirectMax,
		KeyFetchBodySizeMax,
		KeySearchBodySizeMax,
		KeyDebug,
	}
}

func init() {
	Conf.MyType = "TypeConf"

	Vip.SetConfigType("json")
	Vip.SetEnvPrefix(EnvPrefix)
	Vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	Vip.AutomaticEnv()

	Vip.SetDefault(KeyListenAddress, DefaultListenAddress)
	Vip.SetDefault(KeySearXngUrl, DefaultSearXngUrl)
	Vip.SetDefault(KeyRedirectMax, DefaultRedirectMax)
	Vip.SetDefault(KeyFetchBodySizeMax, DefaultFetchBodySizeMax)
	Vip.SetDefault(KeySearchBodySizeMax, DefaultSearchBodySizeMax)
	Vip.SetDefault(KeyDebug, false)

	// keep the logger correct even before the first Load (e.g. init-time errors)
	SetupLog(false)
}

func (t *TypeConf) New() *TypeConf {
	t.Initialized = true
	t.MyType = "TypeConf"
	return t
}

// Effective - configuration values in use, for display
func (t *TypeConf) Effective() TypeConfEffective {
	return TypeConfEffective{
		FileConf:          t.FileConf,
		ListenAddress:     t.ListenAddress,
		SearXngUrl:        t.SearXngUrl,
		RedirectMax:       t.RedirectMax,
		FetchBodySizeMax:  t.FetchBodySizeMax,
		SearchBodySizeMax: t.SearchBodySizeMax,
		Debug:             t.Debug,
		Version:           Version(),
		Source:            t.Source(),
		Env:               t.Env(),
	}
}

// Load - resolve the effective configuration. Precedence, lowest to highest:
// config file, environment variables, command line flags.
//
// A missing config file is not an error, the defaults are used instead.
func (t *TypeConf) Load() error {
	prefix := t.MyType + ".Load"
	t.Err = nil

	t.FileConf = Vip.GetString(KeyFileConf)
	if t.FileConf == "" {
		t.FileConf = DefaultConfigFile
	}
	Vip.SetConfigFile(TildeEnvExpand(t.FileConf))
	if err := Vip.ReadInConfig(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		t.Err = errs.New(prefix, err.Error())
	}

	var values typeConfValues
	if t.Err == nil {
		if t.Err = Vip.Unmarshal(&values); t.Err != nil {
			t.Err = errs.New(prefix, t.Err.Error())
		}
	}

	if t.Err == nil {
		t.ListenAddress = values.ListenAddress
		t.SearXngUrl = values.SearXngUrl
		t.Debug = Vip.GetBool(KeyDebug)
		// setup the logger before validating, so a config error is also reported
		// at the level the user asked for
		SetupLog(t.Debug)
		// read the numeric keys directly; Unmarshal does not resolve them reliably
		t.RedirectMax = Vip.GetInt(KeyRedirectMax)
		t.FetchBodySizeMax = Vip.GetInt64(KeyFetchBodySizeMax)
		t.SearchBodySizeMax = Vip.GetInt64(KeySearchBodySizeMax)
		t.Err = t.validate()
	}
	if t.Err == nil {
		t.Initialized = true
	}
	return t.Err
}

// validate - reject values that would break the tools at run time
func (t *TypeConf) validate() error {
	prefix := t.MyType + ".Load"
	switch {
	case t.RedirectMax < 0:
		return errs.New(prefix, fmt.Sprintf("%s must be >= 0, got %d", KeyRedirectMax, t.RedirectMax))
	case t.FetchBodySizeMax <= 0:
		return errs.New(prefix, fmt.Sprintf("%s must be > 0, got %d", KeyFetchBodySizeMax, t.FetchBodySizeMax))
	case t.SearchBodySizeMax <= 0:
		return errs.New(prefix, fmt.Sprintf("%s must be > 0, got %d", KeySearchBodySizeMax, t.SearchBodySizeMax))
	}
	return nil
}
