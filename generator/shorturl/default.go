package shorturl

import (
	"github.com/micro-services-roadmap/uid-generator-go/generator"
	"github.com/micro-services-roadmap/uid-generator-go/generator/generators"
	"github.com/micro-services-roadmap/uid-generator-go/utilu"
	"math/rand"
)

type DefaultShortUrl struct {
	generator.UidGenerator
}

func New() (*DefaultShortUrl, error) {
	gen6, err := generators.NewDefaultUidGenerator(32, 11, 20, rand.Int63n(2^11))
	if err != nil {
		return nil, err
	} else {
		return &DefaultShortUrl{gen6}, nil
	}
}

func NewV6() (*DefaultShortUrl, error) {
	gen6, err := generators.NewDefaultUidGenerator(28, 3, 3, rand.Int63n(2^3), "2020-10-12")
	if err != nil {
		return nil, err
	} else {
		return &DefaultShortUrl{gen6}, nil
	}
}

func NewV7() (*DefaultShortUrl, error) {
	gen6, err := generators.NewDefaultUidGenerator(30, 6, 5, rand.Int63n(2^6))
	if err != nil {
		return nil, err
	} else {
		return &DefaultShortUrl{gen6}, nil
	}
}

func NewV8() (*DefaultShortUrl, error) {
	gen6, err := generators.NewDefaultUidGenerator(32, 8, 9, rand.Int63n(2^8))
	if err != nil {
		return nil, err
	} else {
		return &DefaultShortUrl{gen6}, nil
	}
}

func (c *DefaultShortUrl) ShortUrl() string {
	return utilu.ToBase62(c.GetUID())
}