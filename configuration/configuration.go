package configuration

import (
	"os"
    "fmt"

	"github.com/spf13/viper"

	ermcmd "go-orchestrator/cmd"
)

const (
	APP_CONFIG_FILE      string = "exo-release"
	APP_CONFIG_FILE_TYPE string = "properties"
	APP_CONFIG_FILE_PATH string = "/.eXo/Release"
)

// LoadConfigFile load the needed configuration file
func LoadConfigFile() {
	// Init the conf app
	// Looking for ~/.eXoR/exor-config.properties
	viper.SetConfigName(APP_CONFIG_FILE)
	viper.SetConfigType(APP_CONFIG_FILE_TYPE)
	viper.AddConfigPath(os.Getenv("HOME") + APP_CONFIG_FILE_PATH)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		ermcmd.PrintError("Fatal error while trying to load config file:", err)
	}
}

// DisplayUserConfiguration Show the loaded configuration 
func DisplayUserConfiguration() {
     fmt.Println("key | value: ")
    for _, k := range viper.AllKeys() {
        fmt.Println(""+k + " = "+ viper.GetString(k))
    }
}
