package plugind

import (
	"fmt"
	"log"
)

func main() {
	err := StartPluginsService()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(StatusAllPluginsService())
	err = StopPluginsService()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(StatusAllPluginsService())
}
