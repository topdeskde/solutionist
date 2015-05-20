package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
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

type GradleConfig struct {
	version                 string
	group                   string
	description             string
	internalProjectName     string
	customerName            string
	projectFullName         string
	tasVersion              string
	testCase                string
	customerReferenceNumber string
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
	if !downloadGradleBuildTemplate() {
		log.Critical("As the download failed there is nothing more to do. This ends now. :(")
	} else {
		patchGradleConfig()

	}
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
	consoleBackendFormatted := logging.NewBackendFormatter(consoleBackend, consoleFormat)
	consoleBackendLeveled := logging.AddModuleLevel(consoleBackendFormatted)
	consoleBackendLeveled.SetLevel(logging.INFO, "")
	if args.debug {
		logging.SetBackend(consoleBackendFormatted)
	} else {
		logging.SetBackend(consoleBackendLeveled)
	}

	if args.logfile {
		file, err := os.Create(args.dir + "/solutionist.log")
		if err != nil {
			log.Error("%v", err)
		} else {
			fileFormat := logging.MustStringFormatter("%{time:15:04:05.000} %{shortfile:20s} %{level: 8s} | %{message}")
			fileBackend := logging.NewLogBackend(file, "", 0)
			fileBackendFormatted := logging.NewBackendFormatter(fileBackend, fileFormat)
			fileBackendLeveled := logging.AddModuleLevel(fileBackendFormatted)
			fileBackendLeveled.SetLevel(logging.INFO, "")
			if args.debug {
				logging.SetBackend(consoleBackendFormatted, fileBackendFormatted)
			} else {
				logging.SetBackend(consoleBackendLeveled, fileBackendLeveled)
			}

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
	checkEnvVar("GRADLE_HOME")
	checkEnvVar("GRADLE_HOME_USER")
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
	log.Info("> Downloading Gradle build template to directory [%s]", args.dir)

	if args.username == "" {
		log.Warning("Username needed:")
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		input, _, err := reader.ReadLine()
		if err != nil {
			log.Critical("Error: %v", err)
		}
		log.Debug("Username provided: %v", string(input))
		log.Debug("Username provided: %v", input)
		log.Debug("Username length: %d", len(input))
		args.username = string(input)

	}

	if args.password == "" {
		log.Warning("Password needed:")
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		input, _, err := reader.ReadLine()
		if err != nil {
			log.Critical("Error: %v", err)
		}
		log.Debug("Password provided: %v", string(input))
		log.Debug("Password provided: %v", input)
		log.Debug("Password length: %d", len(input))
		args.password = string(input)

	}

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

func collectGradleConfig() {

}

func patchGradleConfig() {
	input, err := ioutil.ReadFile(args.dir + "/gradle.build")
	if err != nil {
		log.Critical("%v", err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if !strings.HasPrefix(line, "/") && !strings.HasPrefix(line, "*") {
			if strings.Contains(line, "]") {
				lines[i] = "LOL"
			}
		}
	}

	output := strings.Join(lines, "\n")
	if err = ioutil.WriteFile("myfile", []byte(output), 0777); err != nil {
		log.Critical("%v", err)
	}
}
