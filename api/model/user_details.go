package model

type UserDetails[T any] struct {
	Id                string
	Data              T
	AuthorityPrefix   string
	AuthoritiesValues []string
}

func NewUserDetails[T any](id string, authValues []string, data T, authPrefix string) UserDetails[T] {
	userDetails := UserDetails[T]{
		Id:                id,
		Data:              data,
		AuthorityPrefix:   authPrefix,
		AuthoritiesValues: authValues,
	}
	return userDetails
}
