// Encode image from local path to base64.
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	bytes, err := ioutil.ReadFile("")
	if err != nil {
		fmt.Println("main.ReadFile,", err)
	}
	fmt.Println(base64.URLEncoding.EncodeToString(bytes))
}