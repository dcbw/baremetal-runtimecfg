package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/openshift/baremetal-runtimecfg/pkg/config"
	"github.com/spf13/cobra"
	"net"
)

var (
	displayCmd = &cobra.Command{
		Use: `display [path to kubeconfig]
			It prints the runtime configuration`,
		Short: "Prints the runtime configuration",
		RunE:  runDisplay,
	}
)

func init() {
	displayCmd.Flags().StringP("cluster-config", "c", "", "Path to cluster-config ConfigMap to retrieve ControlPlane info")
	displayCmd.Flags().Bool("verbose", false, "Display extra information about the rendering")
	displayCmd.Flags().IP("api-vip", nil, "DEPRECATED: Virtual IP Address to reach the OpenShift API")
	displayCmd.Flags().IPSlice("api-vips", nil, "Virtual IP Addresses to reach the OpenShift API")
	displayCmd.Flags().IP("ingress-vip", nil, "DEPRECATED: Virtual IP Address to reach the OpenShift Ingress Routers")
	displayCmd.Flags().IPSlice("ingress-vips", nil, "Virtual IP Addresses to reach the OpenShift Ingress Routers")
	displayCmd.Flags().IP("dns-vip", nil, "DEPRECATED: Virtual IP Address to reach an OpenShift node resolving DNS server")
	displayCmd.Flags().Uint16("api-port", 6443, "Port where the OpenShift API listens at")
	displayCmd.Flags().Uint16("lb-port", 9445, "Port where the API HAProxy LB will listen at")
	displayCmd.Flags().Uint16("stat-port", 29445, "Port where the HAProxy stats API will listen at")
	displayCmd.Flags().StringP("resolvconf-path", "r", "/etc/resolv.conf", "Optional path to a resolv.conf file to use to get upstream DNS servers")
	rootCmd.AddCommand(displayCmd)
}

func runDisplay(cmd *cobra.Command, args []string) error {
	kubeCfgPath := "./kubeconfig"
	if len(args) > 0 {
		kubeCfgPath = args[0]
	}

	apiVip, err := cmd.Flags().GetIP("api-vip")
	if err != nil {
		apiVip = nil
	}
	apiVips, err := cmd.Flags().GetIPSlice("api-vips")
	if err != nil {
		apiVips = []net.IP{}
	}
	// If we were passed a VIP using the old interface, coerce it into the list
	// format that the rest of the code now expects.
	if len(apiVips) < 1 && apiVip != nil {
		apiVips = []net.IP{apiVip}
	}
	ingressVip, err := cmd.Flags().GetIP("ingress-vip")
	if err != nil {
		ingressVip = nil
	}
	ingressVips, err := cmd.Flags().GetIPSlice("ingress-vips")
	if err != nil {
		ingressVips = []net.IP{}
	}
	// If we were passed a VIP using the old interface, coerce it into the list
	// format that the rest of the code now expects.
	if len(ingressVips) < 1 && ingressVip != nil {
		ingressVips = []net.IP{ingressVip}
	}
	apiPort, err := cmd.Flags().GetUint16("api-port")
	if err != nil {
		return err
	}
	lbPort, err := cmd.Flags().GetUint16("lb-port")
	if err != nil {
		return err
	}
	statPort, err := cmd.Flags().GetUint16("stat-port")
	if err != nil {
		return err
	}
	clusterConfigPath, err := cmd.Flags().GetString("cluster-config")
	if err != nil {
		return err
	}

	resolveConfPath, err := cmd.Flags().GetString("resolvconf-path")
	if err != nil {
		return err
	}
	config, err := config.GetConfig(kubeCfgPath, clusterConfigPath, resolveConfPath, apiVips, ingressVips, apiPort, lbPort, statPort)
	if err != nil {
		return err
	}

	spew.Dump(config)
	return err
}
