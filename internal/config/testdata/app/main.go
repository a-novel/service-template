// Command app prints the REST read timeout of the default preset. TestAppPresetDefault runs it
// in a process of its own, since the preset resolves at package initialization.
package main

import (
	"fmt"

	"github.com/a-novel/service-template/internal/config"
)

func main() {
	fmt.Println(config.AppPresetDefault.Rest.Timeouts.Read)
}
