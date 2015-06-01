package main

import (
	"os"
)

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
		log.Warning("%16s"+": %s", key, value)
	} else {
		log.Notice("%16s"+": %s", key, value)
	}
}
