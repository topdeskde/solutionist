package main

// TODO:guess internal projectname from customer and project name
// TODO:guess helga repo path from gradle group
// TODO:try to fix environment if possible
// TODO: check if target dir is empty

import (
	"github.com/op/go-logging"
)

const (
	version = "1.0.1"
)

var (
	log    = logging.MustGetLogger("solutionist")
	args   CmdlineArgs
	gradle GradleConfig
	helga  HelgaConfig
)

func main() {
	args = parseCmdline()
	setupLogging()
	showInfo()
	checkEnvironment()
	downloadGradleBuildTemplate()
	setupDefaultGradleConfig()
	collectGradleConfig()
	patchGradleConfig()
	executeCmd("gradle", `-p`+args.dir+``, "wrapper")
	executeCmd("gradle", `-p`+args.dir+``, "init")
	executeCmd("hg", "init", ``+args.dir+``)
	executeCmd("hg", "addremove", ``+args.dir+``)
	executeCmd("hg", "commit", `-m Start a new Gradle project`, ``+args.dir+``)
	setupDefaultHelgaConfig()
	collectHelgaConfig()
	createHelgaRepo()
}

func showInfo() {
	log.Info(`
            ,    _
           /|   | |
         _/_\_  >_<
        .-\-/.   |
       /  | | \_ |
       \ \| |\__(/
       /('---')  |
      / /     \  |
   _.'  \'-'  /  |
   \----/\=-=/   ' Solutionist ` + version)
	log.Info("========================================")
	log.Debug("")
	log.Debug("Using these parameters (use -h for help):")
	log.Debug("%s", args)
}
