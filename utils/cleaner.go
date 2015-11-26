package utils

import (
	"regexp"
)

func CleanHostname(dirtyHostname string) string {
	match := regexp.MustCompile(`\<?(?:http\:\/\/)?([\w\.\-\_]+)\|?`)
	return match.FindStringSubmatch(dirtyHostname)[1]
}