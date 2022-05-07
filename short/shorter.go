package short

import (
	"github.com/segmentio/ksuid"
	"log"
)

type Shorter interface {
	Generate() string
}

func NewShort(isExist func(string) bool) Shorter {
	k := ksuid.New()
	s := Short{
		isExist: isExist,
		k:       k,
	}
	return &s
}

type Short struct {
	isExist func(shortURL string) bool
	k       ksuid.KSUID
}

func (s *Short) Generate() string {
	s.k = ksuid.New()
	randomId := s.k.String()
	log.Println("randomID is " + randomId)
	short := randomId[0:8]
	return short
}
