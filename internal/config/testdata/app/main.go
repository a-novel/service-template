package main

import (
	"fmt"

	"github.com/a-novel/service-template/internal/config"
)

func main() {
	fmt.Println(config.AppPresetDefault.Rest.Timeouts.Read)
}
