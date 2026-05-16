package status

import (
	"encoding/json"
	"fmt"
	"os"
)

func CreateStatFile(port string) (*os.File, error) {
	file, err := os.Create(os.Getenv("PROJECT_ROOT") + "server-stats/" + port + ".json")
	if err != nil {
		fmt.Println("error from CreateStatFile()")
		return nil, err
	}

	return file, nil
}

func UpdateServerStatus(serverStatusFile *os.File) error {
	metrics := Status()

	if err := serverStatusFile.Truncate(0); err != nil {
		fmt.Println("truncate error from UpdateServerStatus()")
		return err
	}

	if _, err := serverStatusFile.Seek(0, 0); err != nil {
		fmt.Println("seel error from UpdateServerStatus()")
		return err
	}

	return json.NewEncoder(serverStatusFile).Encode(metrics)
}

