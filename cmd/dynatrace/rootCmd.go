package dynatrace

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DtTenantURLFlag = "dynatrace-url"
	DtManagedFlag   = "managed"
)

func NewCmdDynatrace() *cobra.Command {
	dtCmd := &cobra.Command{
		Use:               "dynatrace",
		Aliases:           []string{"dt"},
		Short:             "Dynatrace related utilities",
		Args:              cobra.NoArgs,
		DisableAutoGenTag: true,
	}

	dtCmd.PersistentFlags().String(DtTenantURLFlag, "", "Dynatrace tenant URL (overrides OCM label lookup, e.g. https://myserver.example.com/e/environment-id/)")
	dtCmd.PersistentFlags().Bool(DtManagedFlag, false, "Use Dynatrace Managed authentication (tenant-local OAuth instead of sso.dynatrace.com)")
	_ = viper.BindPFlag("dt_tenant_url", dtCmd.PersistentFlags().Lookup(DtTenantURLFlag))
	_ = viper.BindPFlag("dt_is_managed", dtCmd.PersistentFlags().Lookup(DtManagedFlag))

	dtCmd.AddCommand(NewCmdLogs())
	dtCmd.AddCommand(newCmdURL())
	dtCmd.AddCommand(newCmdDashboard())
	dtCmd.AddCommand(NewCmdHCPMustGather())

	return dtCmd
}
