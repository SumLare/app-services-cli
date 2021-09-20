package connect

import (
	"context"
	"errors"

	"github.com/redhat-developer/app-services-cli/internal/build"
	"github.com/redhat-developer/app-services-cli/internal/config"
	"github.com/redhat-developer/app-services-cli/pkg/cluster"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/factory"
	"github.com/redhat-developer/app-services-cli/pkg/connection"
	"github.com/redhat-developer/app-services-cli/pkg/iostreams"
	"github.com/redhat-developer/app-services-cli/pkg/localize"
	"github.com/redhat-developer/app-services-cli/pkg/logging"
	"github.com/spf13/cobra"
)

type Options struct {
	Config     config.IConfig
	Connection func(connectionCfg *connection.Config) (connection.Connection, error)
	Logger     logging.Logger
	IO         *iostreams.IOStreams
	localizer  localize.Localizer

	kubeconfigLocation string
	namespace          string

	offlineAccessToken      string
	forceCreationWithoutAsk bool
	ignoreContext           bool
	serviceID               string
	serviceType             string
}

func NewConnectCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Config:     f.Config,
		Connection: f.Connection,
		Logger:     f.Logger,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     opts.localizer.MustLocalize("cluster.connect.cmd.use"),
		Short:   opts.localizer.MustLocalize("cluster.connect.cmd.shortDescription"),
		Long:    opts.localizer.MustLocalize("cluster.connect.cmd.longDescription"),
		Example: opts.localizer.MustLocalize("cluster.connect.cmd.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.ignoreContext == true && !opts.IO.CanPrompt() {
				return errors.New(opts.localizer.MustLocalize("flag.error.requiredWhenNonInteractive", localize.NewEntry("Flag", "ignore-context")))
			}
			return runConnect(opts)
		},
	}

	cmd.Flags().StringVar(&opts.kubeconfigLocation, "kubeconfig", "", opts.localizer.MustLocalize("cluster.common.flag.kubeconfig.description"))
	cmd.Flags().StringVar(&opts.offlineAccessToken, "token", "", opts.localizer.MustLocalize("cluster.common.flag.offline.token.description", localize.NewEntry("OfflineTokenURL", build.OfflineTokenURL)))
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "", opts.localizer.MustLocalize("cluster.common.flag.namespace.description"))
	cmd.Flags().BoolVarP(&opts.forceCreationWithoutAsk, "yes", "y", false, opts.localizer.MustLocalize("cluster.common.flag.yes.description"))
	cmd.Flags().BoolVar(&opts.ignoreContext, "ignore-context", false, opts.localizer.MustLocalize("cluster.common.flag.ignoreContext.description"))
	cmd.Flags().StringVar(&opts.serviceID, "service-name", "", opts.localizer.MustLocalize("cluster.common.flag.serviceId.description"))
	cmd.Flags().StringVar(&opts.serviceType, "service-type", "", opts.localizer.MustLocalize("cluster.common.flag.serviceName.description"))

	return cmd
}

func runConnect(opts *Options) error {
	conn, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	clusterConn, err := cluster.NewKubernetesClusterConnection(conn, opts.Config, opts.Logger, opts.kubeconfigLocation, opts.IO, opts.localizer)
	if err != nil {
		return err
	}

	arguments := &cluster.ConnectArguments{
		OfflineAccessToken:      opts.offlineAccessToken,
		ForceCreationWithoutAsk: opts.forceCreationWithoutAsk,
		IgnoreContext:           opts.ignoreContext,
		Namespace:               opts.namespace,
		SelectedService:         opts.serviceType,
		SelectedServiceID:       opts.serviceID,
	}

	connectOpts := cluster.Options{
		IO:         opts.IO,
		Logger:     opts.Logger,
		Localizer:  opts.localizer,
		Config:     opts.Config,
		Connection: conn,
	}

	err = clusterConn.Connect(context.Background(), arguments, connectOpts)
	if err != nil {
		return err
	}

	return nil
}
