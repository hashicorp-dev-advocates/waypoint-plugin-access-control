package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type BuildConfig struct {
	Trigger Trigger  `hcl:trigger, optional`
	Access  []Access `hcl:access, optional`
}

type Trigger struct {
	Config      Config `hcl:config, optional`
	Environment string `hcl:environment, optional`
	Approval    string `hcl:approval, optional`
}

type Config struct {
	Platform string   `hcl:platform, optional`
	Event    string   `hcl:event, optional`
	Teams    []string `hcl:teams, optional`
}

type Access struct {
	Database string `hcl:database, optional`
	Cloud    string `hcl:cloud, optional`
	Account  string `hcl:account, optional`
	Role     string `hcl:role, optional`
}

type Builder struct {
	config BuildConfig
}

// Implement Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// Implement ConfigurableNotify
func (b *Builder) ConfigSet(config interface{}) error {
	_, ok := config.(*BuildConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit
		return fmt.Errorf("Expected *BuildConfig as parameter")
	}

	if b.config.Trigger.Environment == "" {
		return fmt.Errorf("environment must be set")
	}

	if b.config.Trigger.Approval != "" &&
		b.config.Trigger.Approval != "manual" &&
		b.config.Trigger.Approval != "automatic" &&
		b.config.Trigger.Approval != "alert" {
		return fmt.Errorf("invalid value set for approval. must be set to `manual`, `automatic`, or `alert`")
	}

	if b.config.Access == nil {
		fmt.Errorf("access for incident response must be configured")
	}

	for i := range b.config.Access {

		if b.config.Access[i].Cloud != "" &&
			b.config.Access[i].Account == "" {
			fmt.Errorf("`account` must be set when configuring `cloud` access")
		}

		if b.config.Access[i].Role == "valid `role` in vault must be set" {

		}
	}
	return nil
}

// Implement Builder
func (b *Builder) BuildFunc() interface{} {
	return b.build
}

func (b *Builder) build(ctx context.Context, ui terminal.UI) (*Binary, error) {
	u := ui.Status()
	defer u.Close()
	u.Update("Setting up access controls")

	if b.config.Trigger.Approval == "" {
		b.config.Trigger.Approval = "manual"
	}

	return &Binary{}, nil
}
