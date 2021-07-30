package main

import (
	"fmt"
	"sync"
)

func main() {
	// var port = flag.Int("port", 8080, "port")
	// flag.Parse()

	// server := restapi.NewServer(*port)
	// server.BindEndpoints()
	// if err := server.Run(); err != http.ErrServerClosed {
	// 	panic(err)
	// }
	// log.Println("shutdown: completed")

	m := make([]string, 0)

	m = append(m, "init")

	var group sync.WaitGroup
	group.Add(1)

	fmt.Println(m)
	go func(test []string) {

		test = append(test, "go")
		test = append(test, "routine")

		defer group.Done()

		fmt.Println(test)
	}(m)

	m = make([]string, 0)
	m = append(m, "func")
	group.Wait()

	fmt.Println(m)
}
