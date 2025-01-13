package generators

import (
	"fmt"
	"github.com/micro-services-roadmap/uid-generator-go/generator"
	"github.com/micro-services-roadmap/uid-generator-go/worker"
	"testing"
)

var gtor generator.UidGenerator

func init() {
	if g, err := NewDefaultUidGenerator(28, 11, 24, worker.CloudflareWorkerId.Instance().NextWorkerId()); err != nil {
		panic(err)
	} else {
		gtor = g
	}

}

func BenchmarkGetUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = gtor.GetUID()
	}
}

func TestDefaultUidGenerator_GetUID(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		id := gtor.GetUID()
		fmt.Println(id)
	}
}

func TestDefaultUidGenerator_ParseUID(t *testing.T) {
	fmt.Println(gtor.ParseUID(1115424893924622336))
	fmt.Println(gtor.ParseUID(1115435579803230208))

}