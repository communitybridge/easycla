package utils

// StringSet contains string set
type StringSet struct {
	set map[string]interface{}
}

// NewStringSet return new StringSet
func NewStringSet() *StringSet {
	return &StringSet{set: make(map[string]interface{})}
}

// Add adds the string to string set
func (ss *StringSet) Add(v string) {
	ss.set[v] = nil
}

// List returns list of strings in StringSet
func (ss *StringSet) List() []string {
	var list []string
	for k := range ss.set {
		list = append(list, k)
	}
	return list
}
