package cmd

import (
	"fmt"
	"os"
	"strings"
	"todolist-api/cmd/http"

	"github.com/spf13/cobra"
)

// AppInfo application info structure
type AppInfo struct {
	AppName        string
	AppVersion     string
	AppCommit      string
	BuildGoVersion string
	BuildArch      string
	BuildDate      string
}

var (
	// meta
	app *AppInfo

	// root command
	rootCmd = &cobra.Command{
		Use:   "todolist-engine-go-sdk",
		Short: "Todolist Engine Go SDK",
		Long:  "Todlist Enginer is Epic!",
	}

	// version sub command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Long:  "Print version information of Todolist",
		Run: func(command *cobra.Command, args []string) {
			infoStr := strings.Builder{}
			infoStr.WriteString(fmt.Sprintf("%s - todlist version info:\n", app.AppName))
			infoStr.WriteString(fmt.Sprintf("Version:\t%s\n", app.AppVersion))
			infoStr.WriteString(fmt.Sprintf("Commit Hash:\t%s\n", app.AppCommit))
			infoStr.WriteString(fmt.Sprintf("Go Version:\t%s\n", app.BuildGoVersion))
			infoStr.WriteString(fmt.Sprintf("Arch:\t\t%s\n", app.BuildArch))
			infoStr.WriteString(fmt.Sprintf("Build:\t\t%s\n", strings.Replace(app.BuildDate, "_", " ", -1)))

			fmt.Println(infoStr.String())
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(http.ServeHTTP())
}

// Execute run root command
func Execute(appInfo *AppInfo) {
	app = appInfo
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// GetAppInfo return application information
func GetAppInfo() *AppInfo {
	return app
}
