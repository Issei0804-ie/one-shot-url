package short

import (
	"github.com/segmentio/ksuid"
	"log"
)

type Shorter interface {
	Generate() string
}

func NewShort() Shorter {
	k := ksuid.New()
	s := Short{
		k: k,
	}
	return &s
}

type Short struct {
	k ksuid.KSUID
}

func (s *Short) Generate() string {
	s.k = ksuid.New()
	randomId := s.k.String()
	log.Println("randomID is " + randomId)
	code := randomId[0:8]
	return code
}
