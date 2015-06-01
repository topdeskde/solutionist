package main

import (
	"github.com/op/go-logging"
	"os"
)

type Hidden string

func (h Hidden) Redacted() interface{} {
	return logging.Redact(string(h))
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
