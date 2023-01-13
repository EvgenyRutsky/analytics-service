package main

import (
	"analytics/app"
	"fmt"
)

func main() {


	application, err := app.New()
	if err != nil {
		fmt.Println(err)
	}

	application.Start()
}
