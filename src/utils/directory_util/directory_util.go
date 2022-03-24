// Package directory_util has all utilities methods related to
// directories of Yako
package directory_util

import (
	"fmt"
	"log"
	"os"
)

// WorkDir checks if the working directory of the node exists. It is created if it does not
// The working directory is located in /usr/<nodeName> and it's used to
// store the uploaded apps
func WorkDir(nodeName string) {
	// Doesn't exist, let's create it
	if _, err := os.Stat("/usr/" + nodeName); os.IsNotExist(err) {
		log.Println(fmt.Sprintf("creating %s's working directory", nodeName))
		if err = os.Mkdir("/usr/"+nodeName, 0775); err != nil {
			log.Println(fmt.Sprintf("creating %s's working directory", nodeName))
			panic(nodeName + "working directory could not be created!")
		}
	}
}
