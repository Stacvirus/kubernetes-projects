package main

import (
	"flag"
	"os"

	webv1alpha1 "github.com/stacvirus/dummy-site-operator/api/v1alpha1"
	"github.com/stacvirus/dummy-site-operator/controllers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		os.Exit(1)
	}

	webv1alpha1.AddToScheme(mgr.GetScheme())

	reconciler := &controllers.DummySiteReconciler{
		Client: mgr.GetClient(),
	}
	reconciler.SetupWithManager(mgr)

	mgr.Start(ctrl.SetupSignalHandler())
}
