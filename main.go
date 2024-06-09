package main

import (
    "flag"
    "os"

    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
    "k8s.io/client-go/kubernetes/scheme"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"

    examplev1 "github.com/toni-marlop/email-operator/api/v1"
    "github.com/toni-marlop/email-operator/controllers"
)

var (
    setupLog = ctrl.Log.WithName("setup")
)

func init() {
    utilruntime.Must(scheme.AddToScheme(scheme.Scheme))
    utilruntime.Must(examplev1.AddToScheme(scheme.Scheme))
}

func main() {
    var metricsAddr string
    var enableLeaderElection bool
    flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
    flag.BoolVar(&enableLeaderElection, "leader-elect", false,
        "Enable leader election for controller manager. "+
            "Enabling this will ensure there is only one active controller manager.")

    flag.Parse()

    ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:                 scheme.Scheme,
        MetricsBindAddress:     metricsAddr,
        Port:                   9443,
        LeaderElection:         enableLeaderElection,
        LeaderElectionID:       "example.com",
        HealthProbeBindAddress: ":8081",
    })
    if err != nil {
        setupLog.Error(err, "unable to start manager")
        os.Exit(1)
    }

    if err = (&controllers.EmailReconciler{
        Client: mgr.GetClient(),
        Scheme: mgr.GetScheme(),
    }).SetupWithManager(mgr); err != nil {
        setupLog.Error(err, "unable to create controller", "controller", "Email")
        os.Exit(1)
    }

    setupLog.Info("starting manager")
    if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
        setupLog.Error(err, "problem running manager")
        os.Exit(1)
    }
}

