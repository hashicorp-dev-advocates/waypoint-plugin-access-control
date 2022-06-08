package release

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// Implement the Destroyer interface
func (rm *ReleaseManager) DestroyFunc() interface{} {
	return rm.destroy
}

func (rm *ReleaseManager) destroy(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	release *Release,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	// Nothing to destroy
	return nil
}
