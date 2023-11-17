package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuxant/voyager-operator/api/v1"
	"github.com/xuxant/voyager-operator/controllers"
	"github.com/xuxant/voyager-operator/pkg/configuration/base/resources"
	"github.com/xuxant/voyager-operator/pkg/log"
	"github.com/xuxant/voyager-operator/version"
	ingressv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	r "runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	metricsHost       = "0.0.0.0"
	metricsPort int32 = 8383
	scheme            = runtime.NewScheme()
	logger            = logf.Log.WithName("cmd")
)

func printInfo() {
	logger.Info(fmt.Sprintf("Version: %s", version.Version))
	logger.Info(fmt.Sprintf("Git Commit: %s", version.GitCommit))
	logger.Info(fmt.Sprintf("Go Version: %s", r.Version()))
	logger.Info(fmt.Sprintf("Go OS/Arch: %s/%s", r.GOOS, r.GOARCH))
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(ingressv1.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string

	isRunningInCluster, err := resources.IsRunningInCluster()
	if err != nil {
		fatal(errors.Wrap(err, "failed to determine if the operator is running in cluster"), true)
	}

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", isRunningInCluster, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
	kubernetesClusterDomain := flag.String("cluster-domain", "cluster.local", "Use custom domain name instead of 'cluster.local'.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	debug := &opts.Development
	log.Debug = *debug
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	printInfo()

	cfg, err := config.GetConfig()
	if err != nil {
		fatal(errors.Wrap(err, "failed to get config"), *debug)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: probeAddr,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "c673222f.tlon.io",
	})
	if err != nil {
		fatal(errors.Wrap(err, "unable to start manager"), *debug)
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fatal(errors.Wrap(err, "failed to setup kubernetes client set"), *debug)
	}
	if resources.IsIngressAPIAvailable(clientSet) {
		logger.Info("Ingress API found: Ingress creation will be performed.")
	}

	if *kubernetesClusterDomain == "" {
		fatal(errors.Wrap(err, "kubernetes cluster domain cannot be empty"), *debug)
	}

	if err = (&controllers.ShipReconciler{
		Client:                  mgr.GetClient(),
		Scheme:                  mgr.GetScheme(),
		ClientSet:               *clientSet,
		Config:                  *cfg,
		KubernetesClusterDomain: *kubernetesClusterDomain,
	}).SetupWithManager(mgr); err != nil {
		fatal(errors.Wrap(err, "unable to create ship controller"), *debug)
	}

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		fatal(errors.Wrap(err, "unable to setup health check"), *debug)
	}

	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		fatal(errors.Wrap(err, "unable to setup ready check"), *debug)
	}

	logger.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fatal(errors.Wrap(err, "problem running manager"), *debug)
	}

}
func fatal(err error, debug bool) {
	if debug {
		logger.Error(nil, fmt.Sprintf("%+v", err))
	} else {
		logger.Error(nil, fmt.Sprintf("%s", err))
	}
}
