package lib

import (
    "strings"
)

func CleanName(name string) string {
    return name
}

func SplitName(name string) []string {
    return strings.Split(name, " ")
}

func MatchNames(name string) bool {
    return false
}
