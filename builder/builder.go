package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type BuildConfig struct {
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

	return &Binary{}, nil
}
