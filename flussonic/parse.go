package flussonic

import (
	"fmt"
	"github.com/mef13/flussonic_exporter/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/url"
)

type Flussonic struct {
	Url            *url.URL
	User           string
	Password       string
	ScrapeInterval string
	InstanceName   string
}

func ParseConfig(v *viper.Viper, key string) ([]*Flussonic, error) {
	type f struct {
		Url            string `mapstructure:"url"`
		User           string `mapstructure:"user"`
		Password       string `mapstructure:"password"`
		ScrapeInterval string `mapstructure:"scrape-interval"`
		InstanceName   string `mapstructure:"instance-name"`
	}

	if v == nil {
		return nil, fmt.Errorf("flussonic configuration not found")
	}

	var confs []f
	err := v.UnmarshalKey(key, &confs)
	if err != nil {
		return nil, err
	}

	var fluss []*Flussonic
	for _, conf := range confs {
		flussUrl, err := url.Parse(conf.Url)
		if err != nil {
			logger.Error("error parsing flussonic url", zap.String("url", conf.Url))
			return nil, err
		}
		if conf.ScrapeInterval == "" {
			conf.ScrapeInterval = "60s"
		}
		if conf.InstanceName == "" {
			conf.InstanceName = flussUrl.Host
		}
		fluss = append(fluss, &Flussonic{
			Url:            flussUrl,
			User:           conf.User,
			Password:       conf.Password,
			ScrapeInterval: conf.ScrapeInterval,
			InstanceName:   conf.InstanceName,
		})
	}
	if len(fluss) == 0 {
		return nil, fmt.Errorf("flussonic configuration not found")
	}
	return fluss, nil
}
