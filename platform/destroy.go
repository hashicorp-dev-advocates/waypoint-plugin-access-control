package platform

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// Implement the Destroyer interface
func (p *Platform) DestroyFunc() interface{} {
	return p.destroy
}

func (p *Platform) destroy(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	deployment *Deployment,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	// Nothing to destroy
	return nil
}
