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

// Length is the length of string set
func (ss *StringSet) Length() int {
	return len(ss.set)
}

// Include tells if the string to string set
func (ss *StringSet) Include(v string) bool {
	_, ok := ss.set[v]
	return ok
}

// List returns list of strings in StringSet
func (ss *StringSet) List() []string {
	var list []string
	for k := range ss.set {
		list = append(list, k)
	}
	return list
}

// Include check is string is present in set or not
func (ss *StringSet) Include(k string) bool {
	_, ok := ss.set[k]
	return ok
}
