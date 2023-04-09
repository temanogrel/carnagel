package internal

import "strings"

func InstanceToIp(instance string) string {
	return strings.Split(instance, ":")[0]
}
