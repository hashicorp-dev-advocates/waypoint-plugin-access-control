package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type Builder struct {
	VaultAddr         string
	VaultToken        string
	VaultNamespace    string
	GrafanaAddr       string
	GrafanaDatasource string
	config            Config
}

func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

func (b *Builder) ConfigSet(config interface{}) error {
	_, ok := config.(*Config)
	if !ok {
		return fmt.Errorf("expected *Config as parameter")
	}

	for _, trigger := range b.config.Trigger {
		switch trigger.Approval {
		case
			"manual",
			"automatic",
			"alert":
		default:
			return fmt.Errorf("invalid value set for approval. must be set to `manual`, `automatic`, or `alert`")
		}

		if trigger.Environment == "" {
			return fmt.Errorf("environment must be set")
		}
	}

	for _, access := range b.config.Access {
		if access.Cloud != "" && access.Account == "" {
			return fmt.Errorf("`account` must be set when configuring `cloud` access")
		}

		if access.Role == "" {
			return fmt.Errorf("valid `role` in vault must be set")
		}
	}

	return nil
}

func (b *Builder) BuildFunc() interface{} {
	return b.build
}

func (b *Builder) build(ctx context.Context, ui terminal.UI, job *component.JobInfo) (*Output, error) {
	u := ui.Status()
	defer u.Close()

	for _, trigger := range b.config.Trigger {
		b.CreateAlert(job.Project, trigger)
	}

	for _, access := range b.config.Access {
		b.CreatePolicy(access)
	}

	return &Output{}, nil
}
