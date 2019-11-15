package comm

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"github.com/spf13/viper"
)

func init() {

}

type Config struct {
	mtx *sync.Mutex
}

func NewConfig() *Config {

	return &Config{mtx: new(sync.Mutex)}
}

//Support JSON YAML TOML HCL
func (this *Config) Open(strFileName string) {

	if this.mtx == nil {
		this.mtx = new(sync.Mutex)
	}
	this.mtx.Lock()
	prefix := "./"
	var name  string
	if strFileName != "" {
		prefix, name = filepath.Split(strFileName)
	}
	name = strings.TrimSuffix(name, filepath.Ext(name))
	viper.SetConfigName(name)
	viper.AddConfigPath(prefix)
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: err=%s, prefix=%s, name=%s \n", err.Error(), prefix, name))
	}
	this.mtx.Unlock()
}

func (this *Config) SetDefault(key string, value interface{}) {
	this.mtx.Lock()
	viper.SetDefault(key, value)
	this.mtx.Unlock()
}

func (this *Config) GetString(key string) (value string){
	this.mtx.Lock()
	value = viper.GetString(key)
	this.mtx.Unlock()
	return
}

// GetBool returns the value associated with the key as a boolean.
func (this *Config) GetBool(key string) (value bool) {
	this.mtx.Lock()
	value =  viper.GetBool(key)
	this.mtx.Unlock()
	return
}

// GetInt returns the value associated with the key as an integer.
func (this *Config) GetInt(key string) (value int) {
	this.mtx.Lock()
	value = viper.GetInt(key)
	this.mtx.Unlock()
	return
}
