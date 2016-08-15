package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mitchellh/colorstring"
	"github.com/spf13/viper"

	ermclient "releases-manager/orchestrator/client"
	ermcmd "releases-manager/orchestrator/cmd"
	ermconf "releases-manager/orchestrator/configuration"
	ermtypes "releases-manager/orchestrator/types"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = ` 
            ███████╗██████╗ ███╗   ███╗
            ██╔════╝██╔══██╗████╗ ████║
            █████╗  ██████╔╝██╔████╔██║
            ██╔══╝  ██╔══██╗██║╚██╔╝██║
            ███████╗██║  ██║██║ ╚═╝ ██║
            ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝
                           
        eXo Releases Manager to orchestre all Releases Process.
        Version: %s
        
`
	// VERSION is the binary version.
	VERSION = "v1.1.0-alpha1"
)

var (
	cmd     string
	catalog ermtypes.Catalog
	label   string
	jiraId  string
	version bool
	// to display help commands
	help bool

	// default values
	defaultCommand = ermcmd.CmdInfo
	defaultID      = "SWF-9999"
	defaultLabel   = ermcmd.LabelAll
)

func init() {

	// parse flags
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&help, "help", false, "print help and exit")
	flag.BoolVar(&help, "h", false, "print help and exit (shorthand)")
	flag.StringVar(&cmd, "command", defaultCommand, "Command to run in each container")
	flag.StringVar(&cmd, "c", defaultCommand, "Command to run in each container (shorthand)")
	flag.StringVar(&label, "label", defaultLabel, "Label to identify a list of projects")
	flag.StringVar(&label, "l", defaultLabel, "Label to identify a list of projects (shorthand)")
	flag.StringVar(&jiraId, "id", defaultID, "JIRA ID")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(colorstring.Color("[blue]"+BANNER), VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()
	// Display help if no flag
	if flag.NArg() > 0 {
		fmt.Println("nb args:", flag.Args())
		ermcmd.UsageAndExit("Pass the command name.", 1)
	}

	// if flags are ok, then check for pre-requisites
	checkPrerequisites()
}

// Check Prerequisites to run this program
// 1- config file
// 2- Docker daemon
func checkPrerequisites() {

	// First load the required config file
	ermconf.LoadConfigFile()

	// Load the Catalog
	var err error
	catalog, err = ermtypes.GetCatalog(viper.GetString(ermclient.CATALOG_BASE_URL) + "/" + defaultID + ".json")
	if err != nil {
		ermcmd.PrintError("Error while donwloading the catalog.", err)
	}
}

func main() {

	if version {
		fmt.Printf("ERM version %s, build %s", VERSION, "eDDdsD2EDSD")
		return
	}

	if help {
		ermcmd.UsageAndExit("", 0)
		return
	}

	log.Printf(colorstring.Color("[yellow]" + "[DEBUG] Running command %s with label %s for JIRA ID %s"), cmd, label, jiraId)
	switch cmd {
    case ermcmd.CmdInfo:
        fmt.Println("TODO; default informations")    
    case ermcmd.CmdDisplayReleasesStatus:
        ermclient.ReadReleaseStatusFromContainer("SWF-3435/task/status")
	case ermcmd.CmdRelease:
		ermclient.ExecAllReleases(catalog, label)
	case ermcmd.CmdDisplayCatalog:
		ermtypes.DisplayCatalog(catalog, label)
	case ermcmd.CmdDisplayUserConfiguration:
		ermconf.DisplayUserConfiguration()
	case ermcmd.CmdDisplayRunningContainers:
		ermclient.ListContainers(false)
	case ermcmd.CmdDisplayAllContainers:
		ermclient.ListContainers(true)
	default:
        fmt.Println(colorstring.Color("[blue]"+BANNER))
		fmt.Printf(colorstring.Color("[red] The command %s is not available. Please try -help to know all commands.\n"), cmd)
	}

	// Hold the execution to look at the events coming
	time.Sleep(5 * time.Second)
}
