package pubsub

import "strings"

type SubjectKey string

const EmptySubject SubjectKey = ""

func NewSubjectKey(elements ...string) SubjectKey {
	return SubjectKey(strings.Join(elements, "."))
}

func (s SubjectKey) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

func (s SubjectKey) Split() []string {
	return strings.Split(string(s), ".")
}

func (s SubjectKey) IsEmpty() bool {
	return s == EmptySubject
}

func (s SubjectKey) String() string {
	return string(s)
}
