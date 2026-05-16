package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App     AppConfig     `yaml:"app"`
	Mine    MineConfig    `yaml:"mine"`
	Factory FactoryConfig `yaml:"factory"`
	Ingot   IngotConfig   `yaml:"ingot"`
	Item    ItemConfig    `yaml:"item"`
	Ore     OreConfig     `yaml:"ore"`
	Storage StorageConfig `yaml:"storage"`
}

type AppConfig struct {
	RunDuration time.Duration `yaml:"run_duration"`
	LogFile     string        `yaml:"log_file"`
}

type MineConfig struct {
	BufferSize int           `yaml:"buffer_size"`
	TickRate   time.Duration `yaml:"tick_rate"`
}

type FactoryConfig struct {
	OutBufferSize int          `yaml:"out_buffer_size"`
	Scaler        ScalerConfig `yaml:"scaler"`
}

type ScalerConfig struct {
	MinWorkers         int           `yaml:"min_workers"`
	MaxWorkers         int           `yaml:"max_workers"`
	ScaleUpThreshold   int           `yaml:"scale_up_threshold"`
	ScaleDownThreshold int           `yaml:"scale_down_threshold"`
	CheckInterval      time.Duration `yaml:"check_interval"`
}

type IngotConfig struct {
	MaxQuality int `yaml:"max_quality"`
}

type ItemConfig struct {
	MaxQuality   int `yaml:"max_quality"`
	QualityBonus int `yaml:"quality_bonus"`
}

type OreConfig struct {
	MaxCapacity  int      `yaml:"max_capacity"`
	NumberOfOres int      `yaml:"number_of_ores"`
	Materials    []string `yaml:"materials"`
}

type StorageConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func ReadConfig(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config

	if err := yaml.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
