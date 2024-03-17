package hitcounter

import "log"

type HitCounter struct {
	name  string
	hits  int
	total int
}

func New(name string) HitCounter {
	return HitCounter{
		name:  name,
		hits:  0,
		total: 0,
	}
}

func (hc *HitCounter) Hit() {
	hc.Inc(true)
}

func (hc *HitCounter) Miss() {
	hc.Inc(false)
}

func (hc *HitCounter) Inc(isHit bool) {
	hc.total++
	if isHit {
		hc.hits++
	}
	log.Printf("[HitCounter/%s] %d/%d", hc.name, hc.hits, hc.total)
}
