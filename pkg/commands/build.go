package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"moul.io/assh/v2/pkg/config"
)

var buildConfigCommand = &cobra.Command{
	Use:   "build",
	Short: "Build .ssh/config",
	RunE:  runBuildConfigCommand,
}

var buildJSONConfigCommand = &cobra.Command{
	Use:   "json",
	Short: "Returns the JSON output",
	RunE:  runBuildJSONConfigCommand,
}

// nolint:gochecknoinits
func init() {
	buildConfigCommand.Flags().BoolP("no-automatic-rewrite", "", false, "Disable automatic ~/.ssh/config file regeneration")
	buildConfigCommand.Flags().BoolP("use-proxy-jump", "", false, "Use ProxyJump instead of ProxyCommand")
	buildConfigCommand.Flags().BoolP("expand", "e", false, "Expand all fields")
	buildConfigCommand.Flags().BoolP("ignore-known-hosts", "", false, "Ignore known-hosts file")
	buildConfigCommand.Flags().BoolP("no-header", "", false, "Suppress header generation")
	_ = viper.BindPFlags(buildConfigCommand.Flags())

	buildJSONConfigCommand.Flags().BoolP("expand", "e", false, "Expand all fields")
	_ = viper.BindPFlags(buildJSONConfigCommand.Flags())
}

func runBuildConfigCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.Open(viper.GetString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}

	if viper.GetBool("expand") {
		for name := range conf.Hosts {
			conf.Hosts[name], err = conf.GetHost(name)
			if err != nil {
				return errors.Wrap(err, "failed to expand hosts")
			}
		}
	}

	if !viper.GetBool("ignore-known-hosts") {
		if conf.KnownHostsFileExists() == nil {
			if err := conf.LoadKnownHosts(); err != nil {
				return errors.Wrap(err, "failed to load known-hosts file")
			}
		}
	}

	if viper.GetBool("no-automatic-rewrite") {
		conf.DisableAutomaticRewrite()
	}

	if viper.GetBool("use-proxy-jump") {
		conf.UseProxyJump()
	}

	if viper.GetBool("no-header") {
		conf.SuppressHeaderGeneration()
	}

	return conf.WriteSSHConfigTo(os.Stdout)
}

func runBuildJSONConfigCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.Open(viper.GetString("config"))
	if err != nil {
		return errors.Wrap(err, "failed to open configuration file")
	}

	if viper.GetBool("expand") {
		for name := range conf.Hosts {
			conf.Hosts[name], err = conf.GetHost(name)
			if err != nil {
				return errors.Wrap(err, "failed to expand hosts")
			}
		}
	}

	s, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	fmt.Println(string(s))
	return nil
}
