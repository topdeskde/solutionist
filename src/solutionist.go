package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nu7hatch/gouuid"
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
	log    = logging.MustGetLogger("solutionist")
	args   CmdlineArgs
	gradle GradleConfig
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
	isXfgProject            string
	testCase                string
	customerReferenceNumber string
	uniqueId                string
	projectType             string
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
		collectGradleConfig()
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

func requestConfigValue(value *string, description string) {
	log.Warning(description)
	log.Info("[%v]", *value)
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _, err := reader.ReadLine()
	if err != nil {
		log.Critical("Error: %v", err)
	}
	log.Debug("Value provided: %v", string(input))
	log.Debug("Value provided: %v", input)
	log.Debug("Value length: %d", len(input))
	if len(input) != 0 {
		*value = string(input)
	}
	log.Notice("Using: %v", *value)

}

func collectGradleConfig() {
	log.Info("> Processing new settings for build.gradle:")

	log.Notice("You can later edit this normally in your editor of choice.")
	log.Notice("The values inside the brackets [] will be used if you enter nothing.")

	uuid4, err := uuid.NewV4()
	if err != nil {
		log.Error("Error while generating UUID:", err)
	}

	gradle = GradleConfig{
		version:                 "1.0.0-SNAPSHOT",
		group:                   "com.topdesk.solution.customer",
		description:             "Tool for customizing icons in the Self Service Desk",
		internalProjectName:     "customer-name_project-name",
		customerName:            "My Customer Name",
		projectFullName:         "My Project Name",
		tasVersion:              "5.5.1",
		isXfgProject:            "false",
		testCase:                "",
		customerReferenceNumber: "",
		uniqueId:                uuid4.String(),
		projectType:             "forms,lookandfeel,labels,reports,modifiedcards,xmlimport,addon,other",
	}

	requestConfigValue(&gradle.version, `
VERSION:
Version of the project, e.g: 1.0.0
Add -SNAPSHOT to indicate it is a work in progress
    `)

	requestConfigValue(&gradle.group, `
GROUP:
One of these depending on the type of your project:
 - com.topdesk.solution.customer (for a TOPdesk client)
 - com.topdesk.solution.addon
 - com.topdesk.solution.prototype
 - com.topdesk.solution.tool (intended for internal use, not limited to consultancy)
 - com.topdesk.solution.lib (a jar not a bespoke zip)
 - com.topdesk.solution.event (like a look & feel for a world cup etc)
 - com.topdesk.solution.product
    `)

	requestConfigValue(&gradle.description, `
DESCRIPTION:
Short description of the project
    `)

	requestConfigValue(&gradle.internalProjectName, `
INTERNALPROJECTNAME:
Used as artifact id for publishing to nexus. Use the format 'customer-name_project-name' if it's a
customer project, otherwise use 'project-name', or 'project-name-x.x' if you release TOPdesk specific        builds (e.g: for an add-on).
    `)

	requestConfigValue(&gradle.customerName, `
CUSTOMERNAME:
Full name of the customer: will end up as part of the ZIP file's name.
    `)

	requestConfigValue(&gradle.projectFullName, `
PROJECTFULLNAME:
Full name of the project: will end up as part of the ZIP file's name.
    `)

	requestConfigValue(&gradle.tasVersion, `
TASVERSION:
The TAS version you want to work on, e.g. 5.4.1
    `)

	requestConfigValue(&gradle.isXfgProject, `
ISXFGPROJECT:
Set this to true if this project uses XFG forms. The zip will be locked automatically. 
This also applies to TOPdesk 5.2+.
    `)

	requestConfigValue(&gradle.testCase, `
TESTCASE:
The test case id associated with this solution (used by TOPdesk's test team).
    `)

	requestConfigValue(&gradle.customerReferenceNumber, `
CUSTOMERREFERENCENUMBER:
The customer reference number of the customer this project is created for. 
You can find this on the customer card in TOPhelp.
    `)

	requestConfigValue(&gradle.uniqueId, `
UNIQUEID:
A unique identifier for your Solution. It can be anything, but it is mandatory when creating a zip. 
SaaS will use this to match old and new versions, and it can also be used by the Portfolio.
It is automatically generated but you can choose to overwrite it.
    `)

	requestConfigValue(&gradle.projectType, `
PROJECTTYPE:
It is not mandatory, but it there will be a warning if it isn’t filled in. 
This makes sure we can categorize Solutions better in the future.

Please provide a comma-separated list of a subset of the following:
forms,lookandfeel,labels,reports,modifiedcards,xmlimport,addon,other
    `)

}

func createNewConfigPart() []string {
    newPart := make([]string, 0)
    append(newPart, "/******************************************")
    append(newPart, "* Generated by Solutionist ` + version + `
    append(newPart, "****************************************** 
    append(newPart, "version     '` + + `
    append(newPart, "group       '` + + `
    append(newPart, "description '` + + `
    append(newPart, "apply plugin: 'solution'
    append(newPart, ""
    append(newPart, "solution {


    )
    
}

func patchGradleConfig() {
	input, err := ioutil.ReadFile(args.dir + "/build.gradle")
	if err != nil {
		log.Critical("%v", err)
	}

	lines := strings.Split(string(input), "\n")

	// find line with 'version', save index
	// find line with 'dependencies', save index
	// comment that part out
	// insert new values before that
	startIndex := 0
	endIndex := 0

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "version") {
			startIndex = i
		}
		if strings.HasPrefix(strings.TrimSpace(line), "dependencies") {
			endIndex = i
		}
	}
    log.Debug("Analyzing build.gradle..")
	log.Debug("'version' found in line %d", startIndex)
	log.Debug("'dependencies' found in line %d", endIndex)

    firstPart := lines[:startIndex]
    oldPart := lines[startIndex:endIndex]
    lastPart := lines[endIndex:]
    
    
    
	/*
	   	for i, line := range lines {
	   		if !strings.HasPrefix(line, "/") && !strings.HasPrefix(line, "*") {
	   			if strings.Contains(line, "]") {
	   				lines[i] = "LOL"
	   			}
	   		}
	   	}

	   	output := strings.Join(lines, "\n")
	   	if err = ioutil.WriteFile(args.dir+"/gradle.build", []byte(output), 0777); err != nil {
	   		log.Critical("%v", err)
	       }
	*/
}
