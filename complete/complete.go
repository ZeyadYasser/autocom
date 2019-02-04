package complete

type Map map[string]interface{}

type AutoCompleter interface {
	Set(string, interface{}) error
	Remove(string) error
	TopN(string, int) (Map, error)
}