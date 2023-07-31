package az

import (
	"fmt"
	"os"
)

func main() {
	f, _ := os.OpenFile("az.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	for _, arg := range os.Args {
		_, _ = f.WriteString(arg + "\n")
	}

	switch os.Args[1] {
	case "version":
		fmt.Println(os.Getenv("cloud_mock_az_version"))
	case "account":
		switch os.Args[2] {
		case "show":
			fmt.Println(os.Getenv("cloud_mock_az_account_show"))
		case "get-access-token":
			fmt.Println(os.Getenv("cloud_mock_az_get_access_token"))
		}
	}
}
