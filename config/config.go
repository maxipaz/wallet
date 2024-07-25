package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// AppConfig struct
type AppConfig struct {
	Blockchain BlockchainConfig
	Contract   ContractConfig
}

// BlockchainConfig struct
type BlockchainConfig struct {
	Address    string `mapstructure:"address"`
	WS         string `mapstructure:"ws"`
	PrivateKey string `mapstructure:"pk"`
	Timeout    string `mapstructure:"timeout"`
	TimeoutIn  time.Duration
}

// ContractConfig struct
type ContractConfig struct {
	Address          string `mapstructure:"address"`
	DefaultWeiFounds int64  `mapstructure:"default_wei_founds"`
}

// environmentPrefix prefix used to avoid environment variable names collisions
const environmentPrefix = "SW"

var (
	// Filename configuration file name.
	Filename string
	// App configuration struct
	App AppConfig

	// environmentVarList list of environment variables read by the app. The name should match with a struct field.
	// The dots will be replaced by underscores, it will be capitalized and the environmentPrefix will be added
	// 		i.e.: blockchain.pk => SW_BLOCKCHAIN_PK
	environmentVarList = []string{
		"blockchain.pk",
	}
)

// Setup bind command flags and environment variables
// The precedence to override a configuration is: flag -> environment variable -> configuration field
func Setup(cmd *cobra.Command, _ []string) error {
	v := viper.New()
	v.SetConfigFile(Filename)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	v.SetEnvPrefix(environmentPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = v.BindPFlag(flag.Name, cmd.Flags().Lookup(flag.Name))
	})
	for _, env := range environmentVarList {
		_ = v.BindPFlag(env, cmd.Flags().Lookup(env))
	}

	if err := v.Unmarshal(&App); err != nil {
		return err
	}

	var err error
	App.Blockchain.TimeoutIn, err = time.ParseDuration(App.Blockchain.Timeout)
	if err != nil {
		return err
	}

	return nil
}
