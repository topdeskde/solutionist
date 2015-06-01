package main

import (
	"bufio"
	"fmt"
	"github.com/bgentry/speakeasy"
	"github.com/nu7hatch/gouuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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

func downloadGradleBuildTemplate() {
	log.Info("")
	log.Info("> Downloading Gradle build template to directory [%s]", args.dir)

	if args.username == "" {
		log.Warning("Username needed:")
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		input, _, err := reader.ReadLine()
		if err != nil {
			log.Critical("Error: %v", err)
			log.Fatal("No reason to go on. This ends now :(")
		}
		log.Debug("Username provided: %v", string(input))
		log.Debug("Username provided: %v", input)
		log.Debug("Username length: %d", len(input))
		args.username = string(input)

	}

	if args.password == "" {
		log.Warning("Password needed:")
		input, err := speakeasy.Ask("> ")
		if err != nil {
			log.Critical("Error: %v", err)
			log.Fatal("No reason to go on. This ends now :(")
		}
		log.Debug("Password provided: %v", Hidden(input))
		log.Debug("Password length: %d", len(input))
		args.password = string(input)

	}

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
			log.Panic("Failed to create build.gradle: %s", err)
		}
		defer file.Close()

		size, err := io.Copy(file, resp.Body)
		if err != nil {
			log.Panic("Failed to write downloaded data to build.gradle: %s", err)
		}

		log.Notice("build.gradle with %v bytes downloaded", size)
	} else {
		log.Critical("Error: %s", resp.Status)
		log.Fatal("No reason to go on. This ends now :(")
	}
}

func setupDefaultGradleConfig() {
	uuid4, err := uuid.NewV4()
	if err != nil {
		log.Error("Error while generating UUID:", err)
	}

	gradle = GradleConfig{
		version:                 "1.0.0-SNAPSHOT",
		group:                   "com.topdesk.solution.customer",
		description:             "Tool for customizing icons in the Self Service Desk",
		customerName:            "Customer Name",
		projectFullName:         "Project Name",
		internalProjectName:     "customer-name_project-name",
		tasVersion:              "5.5.1",
		isXfgProject:            "false",
		testCase:                "",
		customerReferenceNumber: "",
		uniqueId:                uuid4.String(),
		projectType:             "forms,lookandfeel,labels,reports,modifiedcards,xmlimport,addon,other",
	}
}

func collectGradleConfig() {
	log.Info("> Processing new settings for build.gradle:")

	log.Notice("You can later edit this normally in your editor of choice.")
	log.Notice("The values inside the brackets [] will be used if you enter nothing.")

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

	requestConfigValue(&gradle.customerName, `
CUSTOMERNAME:
Full name of the customer: will end up as part of the ZIP file's name.
    `)

	requestConfigValue(&gradle.projectFullName, `
PROJECTFULLNAME:
Full name of the project: will end up as part of the ZIP file's name.
    `)

	requestConfigValue(&gradle.internalProjectName, `
INTERNALPROJECTNAME:
Used as artifact id for publishing to nexus. Use the format 'customer-name_project-name' if it's a
customer project, otherwise use 'project-name', or 'project-name-x.x' if you release TOPdesk specific        builds (e.g: for an add-on).
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
	newPart = append(newPart, "/******************************************")
	newPart = append(newPart, " Generated by Solutionist "+version)
	newPart = append(newPart, " ******************************************/")
	newPart = append(newPart, "version     '"+gradle.version+"'")
	newPart = append(newPart, "group       '"+gradle.group+"'")
	newPart = append(newPart, "description '"+gradle.description+"'")
	newPart = append(newPart, "")
	newPart = append(newPart, "apply plugin: 'solution'")
	newPart = append(newPart, "")
	newPart = append(newPart, "solution {")
	newPart = append(newPart, "    internalProjectName '"+gradle.internalProjectName+"'")
	newPart = append(newPart, "    customerName '"+gradle.customerName+"'")
	newPart = append(newPart, "    projectFullName '"+gradle.projectFullName+"'")
	newPart = append(newPart, "    tasVersion '"+gradle.tasVersion+"'")
	newPart = append(newPart, "    isXfgProject "+gradle.isXfgProject)
	newPart = append(newPart, "    testCase '"+gradle.testCase+"'")
	newPart = append(newPart, "    customerReferenceNumber '"+gradle.customerReferenceNumber+"'")
	newPart = append(newPart, "    uniqueId '"+gradle.uniqueId+"'")
	newPart = append(newPart, "    projectType '"+gradle.projectType+"'")
	newPart = append(newPart, "}")
	newPart = append(newPart, "")
	newPart = append(newPart, "/******************************************")
	newPart = append(newPart, " Commented out by Solutionist "+version)
	newPart = append(newPart, " ******************************************")
	return newPart
}

func patchGradleConfig() {
	input, err := ioutil.ReadFile(args.dir + "/build.gradle")
	if err != nil {
		log.Fatalf("Could not read build.gradle: %s", err)
	}

	lines := strings.Split(string(input), "\n")
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
	newPart := createNewConfigPart()
	oldPart := lines[startIndex:endIndex]
	lastPart := lines[endIndex:]

	newLines := make([]string, 0)
	newLines = append(newLines, firstPart...)
	newLines = append(newLines, newPart...)
	newLines = append(newLines, oldPart...)
	newLines = append(newLines, " ******************************************/", "")
	newLines = append(newLines, lastPart...)

	output := strings.Join(newLines, "\n")
	if err = ioutil.WriteFile(args.dir+"/build.gradle", []byte(output), 0777); err != nil {
		log.Critical("Could not write to build.gradle: %s", err)
	}
}