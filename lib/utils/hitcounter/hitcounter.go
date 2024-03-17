package hitcounter

type HitCounter struct {
	name  string
	hits  int
	total int
	incs  int
}

func New(name string) HitCounter {
	return HitCounter{
		name:  name,
		hits:  0,
		total: 0,
		incs:  0,
	}
}

func (hc *HitCounter) Hit() {
	hc.Inc(true)
}

func (hc *HitCounter) Miss() {
	hc.Inc(false)
}

func (hc *HitCounter) Inc(isHit bool) {
	hc.incs++
	hc.total++
	if isHit {
		hc.hits++
	}
	// log.Printf("[HitCounter/%s] %d/%d", hc.name, hc.hits, hc.total)
}

func (hc *HitCounter) HitN(n int) {
	hc.IncN(true, n)
}

func (hc *HitCounter) MissN(n int) {
	hc.IncN(false, n)
}

func (hc *HitCounter) IncN(isHit bool, n int) {
	hc.incs++
	hc.total += n
	if isHit {
		hc.hits += n
	}
	// log.Printf("[HitCounter/%s] %d/%d", hc.name, hc.hits, hc.total)
}
