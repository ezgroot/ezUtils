package zhttp

import (
	"fmt"
	"net/http"
)

func CloseRsp(response *http.Response) {
	if response == nil {
		return
	}

	if response.Body == nil {
		return
	}

	err := response.Body.Close()
	if err != nil {
		fmt.Printf("close response body error = %s\n", err)
	}
}
