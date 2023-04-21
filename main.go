package main

import (
	"fmt"
	"runtime"
	"time"
	"todolist-api/cmd"
)

var (
	// application metadata
	appName    = "Todolist"
	appVersion = "development"
	appCommit  = "xxxxxxx"
	goVersion  = runtime.Version()
	buildDate  = time.Now().UTC().Format("2006-01-02_15:04:05_UTC")
	buildArch  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

func getAppInfo() *cmd.AppInfo {
	if appVersion == "" {
		appVersion = "0.0.1"
	}

	return &cmd.AppInfo{
		AppName:        appName,
		AppVersion:     appVersion,
		AppCommit:      appCommit,
		BuildGoVersion: goVersion,
		BuildArch:      buildArch,
		BuildDate:      buildDate,
	}
}

// main function
func main() {
	cmd.Execute(getAppInfo())
}
