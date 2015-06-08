# Solutionist

Internal project to automate the steps described here:
https://topwiki.topdesk.com/wiki/Gradle_Solution_Plugin

## Prerequisites
Solutionist is a commandline application which interacts with
 * Nexus
 * Gradle
 * Mercurial
 * SCM-Manager

## How to use

The Solutionist will try to guess certain things and ask about others.
If not specified it will use the current directory as its working directory.
It will use the username of the current user to login to HELGA.
It will use the username of the current user to assign a contact mail address.

This behavior can be overridden with commandline flags:
 * dir - sets the project directory. Use quotes if the path contains blanks
 * username - sets the username for Helga, usually guessed
 * password - sets the password for Helga
 * logfile - creates a logfile in the project directory
 * debug - provides some additional information
 * color - use colors in output

 The 'color' flag requires an ANSI-capable terminal. Have a look at the [cmder].

```
solutionist -dir="d:\my funky project" -username=chuckn -password=iamchucknorris -logfile -debug
```

The preferred way to execute the solutionist is to have it on the %path%. Then you can simply use a terminal
to navigate to the desired target directory and execute it without the 'dir' parameter.

[cmder]: http://gooseberrycreative.com/cmder/ "Cmder"



