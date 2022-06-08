package release

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type ReleaseConfig struct {
}

type ReleaseManager struct {
	config ReleaseConfig
}

func (r *Release) URL() string { return "" }

// Implement Configurable
func (rm *ReleaseManager) Config() (interface{}, error) {
	return &rm.config, nil
}

// Implement ConfigurableNotify
func (rm *ReleaseManager) ConfigSet(config interface{}) error {
	_, ok := config.(*ReleaseConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit
		return fmt.Errorf("Expected *ReleaseConfig as parameter")
	}

	return nil
}

// Implement Builder
func (rm *ReleaseManager) ReleaseFunc() interface{} {
	// return a function which will be called by Waypoint
	return rm.release
}

func (rm *ReleaseManager) release(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
) (*Release, error) {
	u := ui.Status()
	defer u.Close()
	u.Update("Skipped step")

	return &Release{}, nil
}
