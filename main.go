package main

import (
	"context"

	"github.com/mircearem/resourcer/rh"
)

func main() {
	ctx := context.Background()
	rhandler := rh.NewHandler(ctx)

}
