package main

import (
	"flag"
	//fp "path/filepath"
	"github.com/op/go-logging"
	"io"
	"net/http"
	"os"
	user "os/user"
	"strings"
)

const (
	version = "0.0.1"
)

var (
	log  = logging.MustGetLogger("solutionist")
	args CmdlineArgs
)

type CmdlineArgs struct {
	dir      string
	username string
	password string
	logfile  bool
	debug    bool
	//port int
}

type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

func main() {
	args = parseCmdline()
	setupLogging()
	showInfo()
	checkEnvironment()
	//checkEmptyDir()
	downloadGradleBuildTemplate()
}

func parseCmdline() CmdlineArgs {
	defaultUsername := ""
	currentUser, err := user.Current()
	if err == nil {
		usernameParts := strings.Split(currentUser.Username, "\\")
		defaultUsername = usernameParts[len(usernameParts)-1]
	}

	dir := flag.String("dir", ".", "Target directory to create project in; defaults to current directory")
	username := flag.String("username", defaultUsername, "Username used for authentication")
	password := flag.String("password", "", "Password used for authentication")
	logfile := flag.Bool("logfile", false, "Logs output to logile in project directory")
	debug := flag.Bool("debug", false, "Show debug information")
	flag.Parse()

    err = os.MkdirAll(*dir, 0777)
	if err != nil {
		panic(err)
	}
    
	return CmdlineArgs{dir: *dir, username: *username, password: *password, logfile: *logfile, debug: *debug}
}


func setupLogging() {
	consoleFormat := logging.MustStringFormatter("%{color}%{message}%{color:reset}")
	consoleBackend := logging.NewLogBackend(os.Stdout, "", 0)
	consoleBackendFormatter := logging.NewBackendFormatter(consoleBackend, consoleFormat)
	logging.SetBackend(consoleBackendFormatter)

	if args.logfile {
		file, err := os.Create(args.dir +"/solutionist.log")
		if err != nil {
			log.Error("%v", err)
		} else {
			fileFormat := logging.MustStringFormatter("%{time:15:04:05.000} %{shortfile:20s} %{level: 8s} | %{message}")
			fileBackend := logging.NewLogBackend(file, "", 0)
			fileBackendFormatter := logging.NewBackendFormatter(fileBackend, fileFormat)
			logging.SetBackend(consoleBackendFormatter, fileBackendFormatter)
		}
	}
}

func showInfo() {
	log.Info("=====================")
	log.Info("| Solutionist %s |", version)
	log.Info("=====================")
	log.Debug("")
	log.Debug("Using these parameters (use -h for help):")
	log.Debug("%+v", args)
}

func checkEnvironment() {
	log.Info("")
	log.Info("> Checking environment:")
	checkEnvVar("JAVA_HOME")
	checkEnvVar("JAVA_HOME_6")
	checkEnvVar("JAVA_HOME_7")
	checkEnvVar("JAVA_HOME_8")
	log.Debug("TODO: Offer to fix that up if necessary")
}

func checkEnvVar(key string) {
	value := os.Getenv(key)
	if value == "" {
		value = "NOT SET"
		log.Warning(key + ":\t" + value)
	} else {
		log.Notice(key + ":\t" + value)
	}
}

func downloadGradleBuildTemplate() bool {
	log.Info("")
	log.Info("> Downloading Gradle build template:")

	log.Debug("TODO: Check if login and password are present")

	success := false
	url := "http://helga/scm/hg/gradle/solution-plugin/raw-file/tip/setup/template-build.gradle"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(args.username, args.password)
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		file, err := os.Create(args.dir + "/build.gradle")
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()

		size, err := io.Copy(file, resp.Body)
		if err != nil {
			log.Panic(err)
		}

		log.Notice("build.gradle with %v bytes downloaded", size)
		success = true
	} else {
		log.Critical(resp.Status)
	}

	return success
}
