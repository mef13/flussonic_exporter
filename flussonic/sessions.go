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
	"net/http"
	"strings"
	"time"
)

type Sessions struct {
	RequestDuration float64
	Url             string
	Sessions        map[string]*MediaSessions
}

type MediaSessions struct {
	Name         string
	DvrClients   float64
	TotalClients float64
	Types        map[string]float64
}

func (f *Flussonic) GetSessions() (*Sessions, error) {
	client := &http.Client{}
	sessions := Sessions{Sessions: make(map[string]*MediaSessions)}
	sessions.Url = "/flussonic/api/sessions"
	req, err := http.NewRequest("GET", f.Url.String()+sessions.Url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(f.User, f.Password)
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	sessions.RequestDuration = time.Since(startTime).Seconds()
	defer resp.Body.Close()

	type entry struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	type head struct {
		Event    string  `json:"event"`
		Sessions []entry `json:"sessions"`
	}

	var entrys head
	err = json.NewDecoder(resp.Body).Decode(&entrys)
	if err != nil {
		return nil, err
	}
	for _, e := range entrys.Sessions {
		if _, ok := sessions.Sessions[e.Name]; !ok {
			sessions.Sessions[e.Name] = &MediaSessions{
				Name:         e.Name,
				DvrClients:   0,
				TotalClients: 0,
				Types:        make(map[string]float64),
			}
		}
		if _, ok := sessions.Sessions[e.Name].Types[e.Type]; !ok {
			sessions.Sessions[e.Name].Types[e.Type] = 0
		}
		sessions.Sessions[e.Name].TotalClients++
		if strings.Contains(e.Type, "dvr") {
			sessions.Sessions[e.Name].DvrClients++
		}
		sessions.Sessions[e.Name].Types[e.Type]++
	}
	return &sessions, nil
}
