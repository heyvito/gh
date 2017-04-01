package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
	"github.com/victorgama/gh/utils"
)

// Open exposes a command responsible for opening repositories on the default
// browser
var Open = cli.Command{
	Name:      "open",
	Aliases:   []string{"o"},
	Usage:     "Opens a repository on GitHub.com",
	UsageText: "gh open [repository]",
	ArgsUsage: "[repository]",
	Action: func(c *cli.Context) error {
		if len(c.Args()) != 1 {
			return fmt.Errorf("usage: gh o (owner/)[repo]")
		}
		rep := utils.RepoURLFromString(c.Args()[0])
		url := fmt.Sprintf("https://github.com/%s", rep.ToURL())

		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		case "windows":
			err = exec.Command(`C:\Windows\System32\rundll32.exe`, "url.dll,FileProtocolHandler", url).Start()
		default:
			err = fmt.Errorf("unsupported platform '%s'. Please open an issue at https://github.com/victorgama/gh", runtime.GOOS)
		}
		return err
	},
}
