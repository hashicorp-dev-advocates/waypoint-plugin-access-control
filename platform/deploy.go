package platform

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type Trigger struct {
	Config      Config `hcl:"config,block"`
	Environment string `hcl:"environment,optional"`
	Approval    string `hcl:"approval"`
}

type Config struct {
	Platform string   `hcl:"platform,optional"`
	Event    string   `hcl:"event,optional"`
	Teams    []string `hcl:"teams,optional"`
}

type Access struct {
	Database string `hcl:"database,optional"`
	Cloud    string `hcl:"cloud,optional"`
	Account  string `hcl:"account,optional"`
	Role     string `hcl:"role,optional"`
}

type DeployConfig struct {
	Trigger []Trigger `hcl:"trigger,block"`
	Access  []Access  `hcl:"access,block"`
}

type Platform struct {
	config DeployConfig
}

// Config Implement Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

var validate *validator.Validate

// ConfigSet Implement ConfigurableNotify
func (p *Platform) ConfigSet(config interface{}) error {
	_, ok := config.(*DeployConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit
		return fmt.Errorf("Expected *BuildConfig as parameter")
	}

	for _, t := range p.config.Trigger {

		switch t.Approval {
		case "manual":
		case "automatic":
		case "alert":
		default:
			return fmt.Errorf("invalid value set for approval. must be set to `manual`, `automatic`, or `alert`")
		}

		if t.Environment == "" {
			return fmt.Errorf("environment must be set")
		}
	}

	for i, t := range p.config.Access {
		validate = validator.New()

		switch {
		case t.Cloud != "":
			if p.config.Access[i].Account == "" {
				fmt.Errorf("`account` must be set when configuring `cloud` access")
			}

		}

		if p.config.Access[i].Role == "" {
			return fmt.Errorf("valid `role` in vault must be set")
		}
	}
	return nil
}

// Implement Builder
func (p *Platform) DeployFunc() interface{} {
	return p.deploy
}

func (p *Platform) deploy(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	dcr *component.DeclaredResourcesResp,
) (*Deployment, error) {
	u := ui.Status()
	defer u.Close()
	u.Update("Deploy access controls")

	u.Update("Access controls deployed successfully")
	ui.Status().Step("Deploying...", "Vault policies are being deployed")
	p.policyWriter()

	for i := range p.config.Access {
		msg := p.config.Access[i].Role + " vault policy successfully written to vault server"
		ui.Status().Step("Deployed", msg)
	}
	return &Deployment{}, nil
}
