package utils

import "context"

var contextKeyID = interface{}("_id_")

func ContextWithID(id interface{}) context.Context {
	return context.WithValue(context.TODO(), contextKeyID, id)
}
