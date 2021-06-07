package tfstate

import (
	"io"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/tfmigrator/tfmigrator/tfmigrator/log"
)

// Updater updates Terraform State by `terraform state` commands.
type Updater struct {
	Stdout    io.Writer
	Stderr    io.Writer
	DryRun    bool
	Logger    log.Logger
	Terraform *tfexec.Terraform
}

func (updater *Updater) logInfo(msg string) {
	if updater.Logger != nil {
		updater.Logger.Info(msg)
	}
}
