package cmd

import (
	"context"
	"os"

	di "github.com/accuknox/accuknox-cli/install"
	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/hubble"
	ci "github.com/cilium/cilium-cli/install"
	"github.com/rs/zerolog/log"

	ki "github.com/kubearmor/kubearmor-client/install"
	"github.com/spf13/cobra"
)

var (
	uninstallOptions ki.Options
	uparams          = ci.UninstallParameters{Writer: os.Stdout}
)

// uninstallCmd represents the get command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall KubeArmor, Cilium and Discovery-engine from a Kubernetes Cluster",
	Long:  `Uninstall KubeArmor, Cilium and Discovery-engine from a Kubernetes Clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Uninstall Discovery-engine
		diOptions.Namespace = "explorer"
		if err := di.DiscoveryEngineUninstaller(client, diOptions); err != nil {
			return err
		}

		// Uninstall KubeArmor
		if err := ki.K8sUninstaller(client, uninstallOptions); err != nil {
			return err
		}

		// Uninstall Cilium
		uparams.Namespace = namespace

		h := hubble.NewK8sHubble(k8sClient, hubble.Parameters{
			Namespace:            uparams.Namespace,
			HelmValuesSecretName: uparams.HelmValuesSecretName,
			RedactHelmCertKeys:   uparams.RedactHelmCertKeys,
			Writer:               uparams.Writer,
		})
		if err := h.Disable(context.Background()); err != nil {
			return err
		}
		uninstaller := ci.NewK8sUninstaller(k8sClient, uparams)
		if err := uninstaller.Uninstall(context.Background()); err != nil {
			log.Error().Msgf("Unable to uninstall Cilium: %s", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().StringVarP(&uninstallOptions.Namespace, "namespace", "n", "kube-system", "Namespace for resources")
	uninstallCmd.Flags().StringVar(&uparams.HelmValuesSecretName, "helm-values-secret-name", defaults.HelmValuesSecretName, "Secret name to store the auto-generated helm values file. The namespace is the same as where Cilium will be installed")
	uninstallCmd.Flags().BoolVar(&uparams.RedactHelmCertKeys, "redact-helm-certificate-keys", true, "Do not print in the terminal any certificate keys generated by helm. (Certificates will always be stored unredacted in the secret defined by 'helm-values-secret-name')")
	uninstallCmd.Flags().StringVar(&uparams.TestNamespace, "test-namespace", defaults.ConnectivityCheckNamespace, "Namespace to uninstall Cilium tests from")
	uninstallCmd.Flags().BoolVar(&uparams.Wait, "wait", false, "Wait for uninstallation to have completed")
}
