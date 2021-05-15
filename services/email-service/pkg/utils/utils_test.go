package utils

import (
	"fmt"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	got := GetEnvVar("DB_USERNAME")
	fmt.Println(got);
}