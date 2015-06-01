package main

import (
	"github.com/franela/goreq"
	"io/ioutil"
	"strings"
)

// tags are used by reflection
type HelgaConfig struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
	Public      string `json:"public"`
}

func setupDefaultHelgaConfig() {
	/* names
	Customer project:                                customers/[reference-number]_[customer-name]/[project-name]
	Add-on:                                          add-ons/[add-on-name]
	Prototype:                                       prototypes/[prototype-name]
	Tool (used by consultants, i.e.: XFG, XIM):      tools/[tool-project-name]
	Libraries:                                       resources/[internal-project-name]
	Playground/Apekooien:                            sandbox/[username]/[project-name]
	*/
	helga = HelgaConfig{
		Name:        "",
		Type:        "hg",
		Description: gradle.description,
		Contact:     args.username + "@topdesk.com",
		Public:      "true",
	}
}

func collectHelgaConfig() {
	log.Info("> Processing settings for new repo on Helga:")
	log.Notice("The values inside the brackets [] will be used if you enter nothing.")

	requestConfigValue(&helga.Name, `
NAME:
One of these depending on the type of your project:
- customers/[reference-number]_[customer-name]/[project-name]
- add-ons/[add-on-name]
- prototypes/[prototype-name]
- tools/[tool-project-name] (Tool, also used by nondevs, e.g. XFG, XIM)
- resources/[internal-project-name] (Libraries go here)
- events/[internal-project-name]
- products/[internal-project-name]
- sandbox/[username]/[project-name] (Playground/Apekooien)

Suggestions are based on the chosen project group.
    `)
	/*
		GROUP:
		One of these depending on the type of your project:
		- com.topdesk.solution.customer (for a TOPdesk client)
		- com.topdesk.solution.addon
		- com.topdesk.solution.prototype
		- com.topdesk.solution.tool (intended for internal use, not limited to consultancy)
		- com.topdesk.solution.lib (a jar not a bespoke zip)
		- com.topdesk.solution.event (like a look & feel for a world cup etc)
		- com.topdesk.solution.product
	*/
}

func createHelgaRepo() {
	res, err := goreq.Request{
		Method:            "POST",
		Uri:               "http://helga/scm/api/rest/repositories",
		BasicAuthUsername: args.username,
		BasicAuthPassword: args.password,
		ContentType:       "application/json",
		Body:              helga,
	}.Do()
	if err != nil {
		log.Fatalf("Could not create repo on Helga: %s", err)
	} else {
		s, _ := res.Body.ToString()
		if s == "" {
			log.Notice("Repository created at: http://helga/scm/hg/%s", helga.Name)
			linkHelgaRepo()
		} else {
			log.Critical("Something went wrong:\n  %v", s)
		}
	}
}

func linkHelgaRepo() {
	hgrc := make([]string, 0)
	hgrc = append(hgrc, "[paths]")
	hgrc = append(hgrc, "default = http://helga/scm/hg/"+helga.Name)

	output := strings.Join(hgrc, "\n")
	if err := ioutil.WriteFile(args.dir+"/.hg/hgrc", []byte(output), 0777); err != nil {
		log.Critical("Could not write to hgrc: %s", err)
	} else {
		log.Notice(".hg/hgrc created accordingly")
	}
}
