package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Datasource struct {
	UID string `json:"uid"`
}

type Folder struct {
	UID   string `json:"uid"`
	Title string `json:"title"`
}

type AlertRule struct {
	UID string `json:"uid"`
}

func (b *Builder) getDatasource() (*Datasource, error) {
	url := fmt.Sprintf("%s/api/datasources/name/%s", b.GrafanaAddr, b.GrafanaDatasource)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var datasource Datasource
	err = json.NewDecoder(response.Body).Decode(&datasource)
	if err != nil {
		return nil, err
	}

	return &datasource, nil
}

func (b *Builder) createFolder(name string) (*Folder, error) {
	payload := Folder{
		UID:   name,
		Title: name,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/folders", b.GrafanaAddr)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var folder Folder
	err = json.NewDecoder(response.Body).Decode(&folder)
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (b *Builder) createAlertRule(name string, folder string, datasource string, environment string, project string, teams string) (*AlertRule, error) {
	data := fmt.Sprintf(`{
		"orgID": 1,
		"folderUID": "%s",
		"ruleGroup": "%s",
		"title": "%s",
		"condition": "B",
		"data": [
			{
				"refId": "A",
				"queryType": "",
				"relativeTimeRange": {
					"from": 600,
					"to": 0
				},
				"datasourceUid": "%s",
				"model": {
					"editorMode": "builder",
					"expr": "count(envoy_server_uptime{job=\"api-deployment\"})",
					"hide": false,
					"intervalMs": 1000,
					"legendFormat": "__auto",
					"maxDataPoints": 43200,
					"range": true,
					"refId": "A"
				}
			},
			{
				"refId": "B",
				"queryType": "",
				"relativeTimeRange": {
					"from": 0,
					"to": 0
				},
				"datasourceUid": "-100",
				"model": {
					"conditions": [
						{
							"evaluator": {
								"params": [
									3
								],
								"type": "lt"
							},
							"operator": {
								"type": "and"
							},
							"query": {
								"params": [
									"A"
								]
							},
							"reducer": {
								"params": [],
								"type": "last"
							},
							"type": "query"
						}
					],
					"datasource": {
						"type": "__expr__",
						"uid": "-100"
					},
					"hide": false,
					"intervalMs": 1000,
					"maxDataPoints": 43200,
					"refId": "B",
					"type": "classic_conditions"
				}
			}
		],
		"noDataState": "NoData",
		"execErrState": "Alerting",
		"for": 300000000000,
		"labels": {
			"environment": "%s",
			"project": "%s",
			"teams": "%s"
		}
	}`, folder, name, name, datasource, environment, project, teams)

	url := fmt.Sprintf("%s/api/v1/provisioning/alert-rules", b.GrafanaAddr)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var rule AlertRule
	err = json.NewDecoder(response.Body).Decode(&rule)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (b *Builder) CreateAlert(project string, trigger Trigger) error {
	// get the prometheus datasource
	datasource, err := b.getDatasource()
	if err != nil {
		return err
	}

	// create a folder for the app
	folder, err := b.createFolder(project)
	if err != nil {
		return err
	}

	// create alert rule
	_, err = b.createAlertRule(trigger.Config.Event, folder.UID, datasource.UID, trigger.Environment, project, strings.Join(trigger.Config.Teams, ","))
	if err != nil {
		return err
	}

	return nil
}
