package utils

import "net/http"

var remoteCli = &http.Client{}

func CloseIdleHttpCli()  {
	remoteCli.CloseIdleConnections()
}