package ports

import (
	"fmt"
	"net"
	"net/http"
)

func ConnectToPort(port int, handler http.Handler) int {
	for {
		addr := fmt.Sprintf(":%d", port)

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			port++
			continue
		}

		fmt.Println("Running on port", port)

		go func() {
			if err := http.Serve(listener, handler); err != nil {
				panic(err)
			}
		}()

		return port
	}
}

