package config

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Config struct {
	*App `json:"app"`
	*Log `json:"log"`
	*DB	`json:"db"`
	*Redis `json:"redis"`

	loadFile string
	reloadInterval time.Duration
}

var Conf *Config

func init() {
	var err error
	Conf, err = NewConfig("/var/work/go-src/person/wegod/config/config.json", time.Second * 10)
	if err != nil {
		panic(err)
	}

	if err = Conf.load(); err != nil {
		panic(err)
	}

	errCh, doLoop, _ := Conf.Loop()
	go func() {
		for {
			select {
			case err = <- errCh:
				panic(err)
			}
		}
	}()

	go doLoop()
}

func NewConfig(loadConfigFile string, reloadInterval time.Duration) (*Config, error) {
	conf := &Config{
		loadFile: loadConfigFile,
		reloadInterval: reloadInterval,
	}

	if err := conf.load(); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) Loop() (errCh chan error, doLoop func(), cancelLoop func()) {
	errCh = make(chan error)
	cancelLoopSignCh := make(chan struct{})
	doLoop = func() {
		var (
			reloadTicker = time.NewTicker(c.reloadInterval)
			isCanceledLoop bool
		)
		for {
			select {
			case _ = <- reloadTicker.C:
				if err := c.load(); err != nil {
					select {
					case errCh <- err:
					default:
						panic(err)
					}
				}
			case _ = <- cancelLoopSignCh:
				reloadTicker.Stop()
				isCanceledLoop = true
			}

			if isCanceledLoop {
				break
			}
		}
	}
	cancelLoop = func() {
		cancelLoopSignCh <- struct{}{}
	}
	return
}

func (c *Config) load() error {
	buf, err := ioutil.ReadFile(c.loadFile)
	if err != nil {
		return err
	}

	switch filepath.Ext(c.loadFile) {
	case ".json":
		err = json.Unmarshal(buf, c)
		if err != nil {
			return err
		}
	case ".toml":
		_, err = toml.Decode(string(buf), c)
		if err != nil {
			return err
		}
	default:
		return errors.New("Either TOML or JSON is supported")
	}

	return nil
}