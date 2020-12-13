package pubsub

import (
	"fmt"
	"strconv"
	"strings"
)

type SubjectKey string

const EmptySubject SubjectKey = ""

func NewSubjectKey(elements ...string) SubjectKey {
	return SubjectKey(strings.Join(elements, "."))
}

func WithUser(key SubjectKey, userId int) SubjectKey {
	return SubjectKey(fmt.Sprintf("%s.user.%d", key.String(), userId))
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

func (s SubjectKey) MustUserId() int {
	tokens := strings.Split(s.String(), ".")
	for i, token := range tokens {
		if token == "user" && i+1 < len(tokens) {
			userId, err := strconv.Atoi(token)
			if err == nil {
				return userId
			}
			break
		}
	}
	panic("could not parse user's id from key" + s.String())
}
