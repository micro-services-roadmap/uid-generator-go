package generators

import (
	"fmt"
	"github.com/micro-services-roadmap/uid-generator-go/generator"
	"sync"
	"time"
)

// DefaultUidGenerator represents the UID generator
type DefaultUidGenerator struct {
	timeBits      int
	workerBits    int
	seqBits       int
	epochStr      string
	epochSeconds  int64
	BitsAllocator *generator.BitsAllocator
	workerId      int64
	sequence      int64
	lastSecond    int64
	mu            sync.Mutex
}

// NewDefaultUidGenerator creates a new DefaultUidGenerator instance
func NewDefaultUidGenerator(timeBits, workerBits, seqBits int, workerId int64, epochStr ...string) (*DefaultUidGenerator, error) {
	//if timeBits+workerBits+seqBits+1 != generator.TotalBits {
	//	return nil, errors.New("the sum of timeBits, workerBits, and seqBits must be 63")
	//}

	gtor := &DefaultUidGenerator{
		timeBits:      timeBits,
		workerBits:    workerBits,
		seqBits:       seqBits,
		BitsAllocator: generator.NewBitsAllocator(timeBits, workerBits, seqBits),
		workerId:      workerId,
	}

	if len(epochStr) == 0 {
		gtor.epochStr = generator.EpochStr
		dt, _ := time.Parse(generator.EpochStrFormat, generator.EpochStr)
		gtor.epochSeconds = dt.Unix()
		return gtor, nil
	}

	if parse, err := time.Parse(generator.EpochStrFormat, epochStr[0]); err != nil {
		gtor.epochStr = generator.EpochStr
		dt, _ := time.Parse(generator.EpochStrFormat, generator.EpochStr)
		gtor.epochSeconds = dt.Unix()
	} else {
		gtor.epochStr = epochStr[0]
		gtor.epochSeconds = parse.Unix()
	}

	return gtor, nil
}

// GetUID generates a unique ID
func (g *DefaultUidGenerator) GetUID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i := 0; i < 10_000; i++ {
		if id, err := g.nextId(); err == nil {
			return id
		}
	}

	panic("UID generation failed")
}

// nextId generates the next UID
func (g *DefaultUidGenerator) nextId() (int64, error) {
	currentSecond, err := g.getCurrentSecond()
	if err != nil {
		return 0, err
	}

	// Handle clock rollback
	if currentSecond < g.lastSecond {
		refusedSeconds := g.lastSecond - currentSecond
		return 0, fmt.Errorf("clock moved backwards. Refusing for %d seconds", refusedSeconds)
	}

	// Increase sequence at the same second
	if currentSecond == g.lastSecond {
		g.sequence = (g.sequence + 1) & g.BitsAllocator.GetMaxSequence()
		// Exceed sequence max, wait for the next second
		if g.sequence == 0 {
			currentSecond = g.getNextSecond(g.lastSecond)
		}
	} else {
		// Reset sequence if it's a new second
		g.sequence = 0
	}

	g.lastSecond = currentSecond

	// Allocate the bits for UID
	return g.BitsAllocator.Allocate(currentSecond-g.epochSeconds, g.workerId, g.sequence), nil
}

// getCurrentSecond gets the current second
func (g *DefaultUidGenerator) getCurrentSecond() (int64, error) {
	currentSecond := time.Now().Unix()
	if currentSecond-g.epochSeconds > g.BitsAllocator.GetMaxDeltaSeconds() {
		return 0, fmt.Errorf("timestamp bits are exhausted. Refusing UID generation")
	}
	return currentSecond, nil
}

// getNextSecond waits for the next second if the current second is exhausted
func (g *DefaultUidGenerator) getNextSecond(lastTimestamp int64) int64 {
	for {
		timestamp := time.Now().Unix()
		if timestamp > lastTimestamp {
			return timestamp
		}
	}
}

// ParseUID parses a UID and returns its components as a string
func (g *DefaultUidGenerator) ParseUID(uid int64) string {
	totalBits := generator.TotalBits
	signBits := g.BitsAllocator.GetSignBits()
	timestampBits := g.BitsAllocator.GetTimestampBits()
	workerIdBits := g.BitsAllocator.GetWorkerIdBits()
	sequenceBits := g.BitsAllocator.GetSequenceBits()

	// Parse UID
	sequence := (uid << uint(totalBits-sequenceBits)) >> uint(totalBits-sequenceBits)
	workerId := (uid << uint(timestampBits+signBits)) >> uint(totalBits-workerIdBits)
	deltaSeconds := uid >> uint(workerIdBits+sequenceBits)

	// Format time from epoch
	thatTime := time.Unix(g.epochSeconds+deltaSeconds, 0)
	return fmt.Sprintf("{\"UID\":\"%d\",\"timestamp\":\"%s\",\"workerId\":\"%d\",\"sequence\":\"%d\"}",
		uid, thatTime.Format("2006-01-02 15:04:05"), workerId, sequence)
}
