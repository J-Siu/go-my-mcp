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

	// Environment variables are prefixed with this, followed by the upper-cased key name
	EnvPrefix = "GO_MY_MCP"
)

const (
	keyDebug         = "Debug"
	keyFileConf      = "FileConf"
	keyListenAddress = "ListenAddress"
	keySearXngUrl    = "SearXngUrl"
)

// Effective configuration, populated by [Conf.Load]
var Conf TypeConf

// Viper instance backing [Conf]. Cobra flags are bound to it from the `cmd` package.
var Vip = viper.New()

type TypeConf struct {
	basestruct.Embedded
	FileConf      string `json:"FileConf"`
	ListenAddress string `json:"ListenAddress"`
	SearXngUrl    string `json:"SearXngUrl"`
}

// Values read from config file, environment and flags
type typeConfValues struct {
	ListenAddress string
	SearXngUrl    string
}

// TypeConfEffective - configuration as shown by the `config` sub command
type TypeConfEffective struct {
	FileConf      string      `json:"FileConf"`
	ListenAddress string      `json:"ListenAddress"`
	SearXngUrl    string      `json:"SearXngUrl"`
	Debug         bool        `json:"Debug"`
	Version       string      `json:"Version"`
	Env           TypeConfEnv `json:"Env"`
}

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

// Env var name for a viper key, e.g. `ListenAddress` -> `GO_MY_MCP_LISTENADDRESS`
func envName(key string) string {
	return EnvPrefix + "_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
}

// Viper keys known to the configuration
func allKeys() []string {
	return []string{keyFileConf, keyListenAddress, keySearXngUrl, keyDebug}
}

func init() {
	Conf.MyType = "TypeConf"

	Vip.SetConfigType("json")
	Vip.SetEnvPrefix(EnvPrefix)
	Vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	Vip.AutomaticEnv()

	Vip.SetDefault(keyListenAddress, DefaultListenAddress)
	Vip.SetDefault(keySearXngUrl, DefaultSearXngUrl)
	Vip.SetDefault(keyDebug, false)
}

func (t *TypeConf) New() *TypeConf {
	t.Initialized = true
	t.MyType = "TypeConf"
	return t
}

// Effective - configuration values in use, for display
func (t *TypeConf) Effective() TypeConfEffective {
	return TypeConfEffective{
		FileConf:      t.FileConf,
		ListenAddress: t.ListenAddress,
		SearXngUrl:    t.SearXngUrl,
		Debug:         t.Debug,
		Version:       Version(),
		Env:           t.Env(),
	}
}

// Load - resolve the effective configuration. Precedence, lowest to highest:
// config file, environment variables, command line flags.
//
// A missing config file is not an error, the defaults are used instead.
func (t *TypeConf) Load() error {
	prefix := t.MyType + ".Load"
	t.Err = nil

	t.FileConf = Vip.GetString(keyFileConf)
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
		t.Debug = Vip.GetBool(keyDebug)
		t.Initialized = true
	}
	return t.Err
}
