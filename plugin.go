package main

import (
	"code.cloudfoundry.org/cli/cf/flags"
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

type cfDeleteWrapper struct{}

// Flags that is used by the cf delete wrapper
type WrapperFlags struct {
	AppName string
	Force   bool
}

// Manifest struct, we only care about the app name
type Manifest struct {
	Applications []struct {
		Name      string `yaml:"name"`
	} `yaml:"applications"`
}

// Subcommand names
const (
	multiAppDeleteCmd         = "delete-multi-apps"
	deleteAppUsingManifestcmd = "delete-app-using-manifest"
)

// Initialize a new flags wrapper
var cmdOptions = new(WrapperFlags)

// Run is what is executed by the Cloud Foundry CLI
func (c *cfDeleteWrapper) Run(cliConnection plugin.CliConnection, args []string) {
	switch args[0] {
	case multiAppDeleteCmd:
		c.buildMultiAppDeleteArguments(args, true)
		c.MultiAppDelete(cliConnection)
	case deleteAppUsingManifestcmd:
		c.buildMultiAppDeleteArguments(args, false)
		c.DeleteAppUsingManifest(cliConnection)
	}
}

// GetMetadata provides the Cloud Foundry CLI with metadata on how to use this plugin
func (c *cfDeleteWrapper) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "CfDeleteWrapper",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     multiAppDeleteCmd,
				Alias:    "dma",
				HelpText: "Delete multiple apps via a single command",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s <APP1>,<APP2>,....,<APPn>", multiAppDeleteCmd),
					Options: map[string]string{
						"app": "-a, list of apps to be deleted",
						"force": "-f, no need to prompt for confirmation",
					},
				},
			},
			{
				Name:     deleteAppUsingManifestcmd,
				Alias:    "daum",
				HelpText: "Detect the apps name from manifest and delete it",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s", deleteAppUsingManifestcmd),
					Options: map[string]string{
						"force": "-f, no need to prompt for confirmation",
					},
				},
			},
		},
	}
}

// Build a list of argument flags to be used by this wrapper
func (c *cfDeleteWrapper) buildMultiAppDeleteArguments(args []string, mandatory bool) {
	fc := flags.New()
	fc.NewBoolFlag("force", "f", "don't prompt for confirmation")
	fc.NewStringFlag("app", "a", "list of apps to delete")
	err := fc.Parse(args[1:]...)
	if err != nil {
		handleError(fmt.Sprintf("%v", err), true)
	}
	if fc.IsSet("f") {
		cmdOptions.Force = true
	}
	if fc.IsSet("app") {
		cmdOptions.AppName = fc.String("app")
	}
	if cmdOptions.AppName == "" && mandatory {
		handleError("Mandatory argument app name is missing", true)
	}
}

// Program that helps to delete multiple apps via single command
func (c *cfDeleteWrapper) MultiAppDelete(cli plugin.CliConnection) {

	// Get the list of apps to be delete
	apps := spiltString(cmdOptions.AppName)

	// Check with user if they want to continue
	if !cmdOptions.Force {
		yesOrNoConfirmation(cmdOptions.AppName)
	}

	// Get the app and delete the app
	checkDeleteApp(cli, apps, false)

}

// Reads the manifest and delete the app
func (c *cfDeleteWrapper) DeleteAppUsingManifest(cli plugin.CliConnection) {

	// Manifest
	m := new(Manifest)

	// Check if the manifest file exists and get the content
	output := readManifest()

	// Get the app name from the manifest
	err := yaml.Unmarshal(output, &m)
	if err != nil {
		handleError(fmt.Sprintf("Encounteed error when unmarshaling the manifest ym, err: %v", err), true)
	}

	// Check & Delete the app name
	if len(m.Applications) > 0 && m.Applications[0].Name != "" {

		// Extract app name
		var apps []string
		for _, a := range m.Applications {
			apps = append(apps, a.Name)
		}

		// Check with user if they want to continue
		if !cmdOptions.Force {
			yesOrNoConfirmation(strings.Join(apps, ","))
		}

		// Get the app and delete the app
		checkDeleteApp(cli, apps, false)

	} else {
		handleError("Unable to find any information from manifest", true)
	}
}

// Check & Delete the app
func checkDeleteApp(cli plugin.CliConnection, apps []string, exit bool) {

	// loop through the app list and delete the app
	for _, app := range apps {

		// Get the app info
		appInfo, _ := cli.GetApp(app)

		// App not found
		if appInfo.Guid == "" {
			handleError(fmt.Sprintf("ERROR: App \"%s\" not found on the current org & space, continuing with remaining apps...", app), exit)
		} else {
			// Run the curl command to delete the app
			output, err  := cli.CliCommandWithoutTerminalOutput("curl", "-X", "DELETE", fmt.Sprintf("/v2/apps/%s", appInfo.Guid))
			if err != nil {
				handleError(fmt.Sprintf("ERROR: Received error when deleting the app \"%s\", err: %v", appInfo.Name, err), exit)
				fmt.Println("continuing with other apps, if there is any...")
			}

			// If we found any message, that means the CLI wants to tell us something is wrong
			if len(output) > 1 {
				handleError(fmt.Sprintf("ERROR: Something went wrong when deleting the app \"%s\", message: \n%v", appInfo.Name, output), exit)
				fmt.Println("continuing with other apps, if there is any...")
			} else {
				fmt.Println("\nSuccessfully deleted the app \"" + appInfo.Name + "\"")
			}
		}
	}
}

// Initialize and start the plugin
func main() {
	plugin.Start(new(cfDeleteWrapper))
}
