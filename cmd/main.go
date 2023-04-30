package main

import (
	"context"
	"flag"

	manager "github.com/ilbarlo/flavourGenerator/pkg/manager"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	ctx := context.Background()

	cl, err := manager.GetKClient(ctx)
	utilruntime.Must(err)

	klog.Fatalf("unable to start the server: %s", manager.SetupRouterAndServeHTTP(ctx, cl))
}
