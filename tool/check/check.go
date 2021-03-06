// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package check

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

//TODO add test
const (
	validIPv4AddressRegex = `(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])`
	validIPv6AddressRegex = `((([0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){6}:[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){5}:([0-9A-Fa-f]{1,4}:)?[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){4}:([0-9A-Fa-f]{1,4}:){0,2}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){3}:([0-9A-Fa-f]{1,4}:){0,3}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){2}:([0-9A-Fa-f]{1,4}:){0,4}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){6}((b((25[0-5])|(1d{2})|(2[0-4]d)|(d{1,2}))b).){3}(b((25[0-5])|(1d{2})|(2[0-4]d)|(d{1,2}))b))|(([0-9A-Fa-f]{1,4}:){0,5}:((b((25[0-5])|(1d{2})|(2[0-4]d)|(d{1,2}))b).){3}(b((25[0-5])|(1d{2})|(2[0-4]d)|(d{1,2}))b))|(::([0-9A-Fa-f]{1,4}:){0,5}((b((25[0-5])|(1d{2})|(2[0-4]d)|(d{1,2}))b).){3}(b((25[0-5])|(1d{2})|(2[0-4]d)|(d{1,2}))b))|([0-9A-Fa-f]{1,4}::([0-9A-Fa-f]{1,4}:){0,5}[0-9A-Fa-f]{1,4})|(::([0-9A-Fa-f]{1,4}:){0,6}[0-9A-Fa-f]{1,4})|(([0-9A-Fa-f]{1,4}:){1,7}:))`
	validHostnameRegex    = `(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])`
)

//IsIP validate ip args format
func IsIP(arg string) bool {
	rIPv4 := regexp.MustCompile(`^` + validIPv4AddressRegex + `$`)
	rIPv6 := regexp.MustCompile(`^` + validIPv6AddressRegex + `$`)
	log.Debug().Str("arg", arg).Bool("isIPv4", rIPv4.MatchString(arg)).Bool("isIPv6", rIPv6.MatchString(arg)).Msg("parsing arg")
	return rIPv4.MatchString(arg) || rIPv6.MatchString(arg)
}

//IsHost validate ip args format
func IsHost(arg string) bool {
	rHost := regexp.MustCompile(`^` + validHostnameRegex + `$`)
	log.Debug().Str("arg", arg).Bool("IsHost", rHost.MatchString(arg)).Msg("parsing arg")
	return rHost.MatchString(arg)
}

//IsValidClientArg validate cient args format
func IsValidClientArg(arg string) bool {
	var host, port string
	if !strings.Contains(arg, ":") {
		//By default use port 8080
		host = arg
		port = "8080"
	} else {
		tmp := strings.SplitN(arg, ":", 2)
		host = tmp[0]
		port = tmp[1]
	}
	log.Debug().Str("arg", arg).Bool("IsIP", IsIP(host)).Bool("IsHost", IsHost(host)).Msg("parsing arg")
	_, err := strconv.Atoi(port)
	return err == nil && (IsIP(host) || IsHost(host))
}

func IsValidFileArg(arg string) bool {
	/*
		file, err := os.Stat(arg)
		if err != nil {
			return false //Failed to open file
		}
		if !file.Mode().IsRegular() {
			return false //Is not a file
		}
	*/
	b, err := ioutil.ReadFile(arg)
	if err != nil {
		return false
	}
	//check whether s contains md comment with json //TODO regex
	return strings.Contains(string(b), "[//]: # ({")
}

func IsValidFolderArg(arg string) bool {
	_, err := ioutil.ReadFile(filepath.Join(arg, "index.md"))
	if err != nil {
		return false
	}
	return true
}

// AskFor ask for validation of a action (maybe if destructive ?)
func AskFor(action string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(action + "? y/n : ")
	text, _ := reader.ReadString('\n')
	for true {
		switch strings.TrimSpace(text) {
		case "y":
			return true
		case "n":
			return false
		}
		fmt.Print(action + "? y/n : ")
		text, _ = reader.ReadString('\n')
	}
	return false
}
