package config

import (
	"github.com/TNK-Studio/lazykube/pkg/utils"
	"github.com/jroimartin/gocui"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
)

var (
	// Todo: check error
	HomePath, _      = utils.Home()
	LazykubeHomePath = path.Join(HomePath, ".lazykube/")

	Conf = &Config{}

	DefaultConfig = &Config{
		GuiConfig: &GuiConfig{
			Highlight:  true,
			Cursor:     false,
			FgColor:    gocui.ColorWhite,
			SelFgColor: gocui.ColorGreen,
			Mouse:      true,
			InputEsc:   true,
		},
		LogConfig: &LogConfig{
			Path:  path.Join(LazykubeHomePath, "log/"),
			Level: logrus.InfoLevel,
		},
		UserConfig: &UserConfig{
			CustomResourcePanels: []string{},
			History: &History{
				ImageHistory:   []string{},
				CommandHistory: []string{},
			},
		},
	}
)

func init() {
	Read()
}

func Read() {
	if !utils.FileExited(LazykubeHomePath) {
		if err := os.MkdirAll(LazykubeHomePath, 0755); err != nil {
			panic(err)
		}
	}

	if err := Conf.ReadFrom(LazykubeHomePath, "config.yaml"); err != nil {
		*Conf = *DefaultConfig
		Save()
	}
}

func Save() {
	if err := Conf.SaveTo(LazykubeHomePath, "config.yaml"); err != nil {
		panic(err)
	}
}

type Config struct {
	GuiConfig  *GuiConfig  `yaml:"gui_config"`
	LogConfig  *LogConfig  `yaml:"log_config"`
	UserConfig *UserConfig `yaml:"user_config"`
}

// ReadFrom read config
func (c *Config) ReadFrom(filePath, fileName string) error {
	configFile, err := ioutil.ReadFile(path.Join(filePath, fileName))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configFile, Conf)
	if err != nil {
		return err
	}

	return nil
}

// SaveTo save config
func (c *Config) SaveTo(filePath, fileName string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(filePath, fileName), bytes, 0755); err != nil {
		return err
	}
	return nil
}
