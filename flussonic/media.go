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
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Media struct {
	RequestDuration float64
	Url             string
	Streams         map[string]*Stream
}

type Stream struct {
	Name    string  `mapstructure:"name"`
	Stats   Stats   `mapstructure:"stats"`
	Options Options `mapstructure:"options"`
}

type Stats struct {
	Bitrate           float64 `mapstructure:"bitrate"`
	Alive             bool    `mapstructure:"alive"`
	ClientCount       float64 `mapstructure:"client_count"`
	DvrEnabled        bool    `mapstructure:"dvr_enabled"`
	InputErrorRate    float64 `mapstructure:"input_error_rate"`
	Lifetime          float64 `mapstructure:"lifetime"`
	RetryCount        float64 `mapstructure:"retry_count"`
	RunningTranscoder bool    `mapstructure:"running_transcoder"`
}

type Options struct {
	Disabled bool   `mapstructure:"disabled"`
	Title    string `mapstructure:"title"`
	Comment  string `mapstructure:"comment"`
}

func (f *Flussonic) GetMedia() (*Media, error) {
	client := &http.Client{}
	media := Media{Streams: make(map[string]*Stream)}
	media.Url = "/flussonic/api/media"
	req, err := http.NewRequest("GET", f.Url.String()+media.Url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(f.User, f.Password)
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	media.RequestDuration = time.Since(startTime).Seconds()
	defer resp.Body.Close()

	type entry struct {
		Entry string      `json:"entry"`
		Value interface{} `json:"value"`
	}
	var entrys []entry
	err = json.NewDecoder(resp.Body).Decode(&entrys)
	if err != nil {
		return nil, err
	}
	for _, e := range entrys {
		if e.Entry == "stream" {
			var stream Stream
			err = mapstructure.Decode(e.Value, &stream)
			if err != nil {
				return nil, err
			}
			media.Streams[stream.Name] = &stream
		}
	}

	return &media, nil
}
