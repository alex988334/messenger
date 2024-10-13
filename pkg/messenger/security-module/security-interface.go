package security_module

type ISecurity interface {
	CheckRequest(request *map[string]interface{}) bool
}
