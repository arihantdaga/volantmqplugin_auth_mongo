package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/VolantMQ/vlapi/vlauth"
	"github.com/VolantMQ/vlapi/vlplugin"
	"gopkg.in/yaml.v2"
)

type authMongoPlugin struct {
	vlplugin.Descriptor
}

// TODO: Add unmarshaller
type config struct {
	MongodbBaseURI string `mapstructure:"mongoURI,omitempty" yaml:"mongoURI,omitempty" json:"mongoURI,omitempty" default:""`
	DatabaseName   string `mapstructure:"database,omitempty" yaml:"database,omitempty" json:"database,omitempty" default:""`
	CollectionName string `mapstructure:"collection,omitempty" yaml:"collection,omitempty" json:"collection,omitempty" default:""`
}

type authImpl struct {
	*vlplugin.SysParams
	*authProvider
}

var _ vlplugin.Plugin = (*authMongoPlugin)(nil)
var _ vlplugin.Info = (*authMongoPlugin)(nil)
var _ vlauth.IFace = (*authImpl)(nil)

// Plugin symbol
var Plugin authMongoPlugin
var version string

func init() {
	Plugin.V = version
	Plugin.N = "mongo"
	Plugin.T = "auth"
}

func (pl *authMongoPlugin) Load(c interface{}, params *vlplugin.SysParams) (pla interface{}, err error) {
	p := &authImpl{
		SysParams: params,
	}
	var cfg config
	decodeIFace := func() error {
		var data []byte
		var e error
		if data, e = yaml.Marshal(c); e != nil {
			e = errors.New(Plugin.T + "." + Plugin.N + ": " + e.Error())
			return e
		}

		if e = yaml.Unmarshal(data, &cfg); e != nil {
			e = errors.New(Plugin.T + "." + Plugin.N + ": " + e.Error())
			return e
		}

		return e
	}

	switch t := c.(type) {
	case map[string]interface{}:
		if err = decodeIFace(); err != nil {
			return nil, err
		}
	case map[interface{}]interface{}:
		if err = decodeIFace(); err != nil {
			return nil, err
		}
	case []byte:
		if err = yaml.Unmarshal(t, &cfg); err != nil {
			err = errors.New(Plugin.T + "." + Plugin.N + ": " + err.Error())
			return nil, err
		}
	default:
		err = fmt.Errorf("%s.%s: invalid config type %s", Plugin.T, Plugin.N, reflect.TypeOf(c).String())
		return nil, err
	}
	p.authProvider = &authProvider{
		cfg: cfg,
	}
	if err = p.authProvider.Init(); err != nil {
		return
	}
	return p, nil

}

func (pl *authMongoPlugin) Info() vlplugin.Info {
	return pl
}

func main() {
	panic("this is a plugin, build it as with -buildmode=plugin")
}
