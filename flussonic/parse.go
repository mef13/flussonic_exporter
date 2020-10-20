/*
 *    Copyright 2020 Yury Makarov
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

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
