package short

import "github.com/segmentio/ksuid"

type Shorter interface {
	Generate() string
}

func NewShort(isExist func(string) bool) Shorter {
	k := ksuid.New()
	s := Short{
		isExist: isExist,
		k:       &k,
	}
	return s
}

type Short struct {
	isExist func(shortURL string) bool
	k       *ksuid.KSUID
}

func (s Short) Generate() string {
	randomId := s.k.Next()
	short := string(randomId[0:8])
	return short
}
