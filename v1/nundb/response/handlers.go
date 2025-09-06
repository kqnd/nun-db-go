package response

type Handler struct {
	command  string
	args     string
	watchers *map[string][]func(interface{})
	pendings *map[string]chan interface{}
}
