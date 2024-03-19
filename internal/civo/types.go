package civo

import (
	"context"

	"github.com/civo/civogo"
)

type CivoConfiguration struct {
	Client  *civogo.Client
	Context context.Context
}

type CivoCmdOptions struct {
	Nuke   bool
	Region string
}
