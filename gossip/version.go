package gossip

type Version struct {
	Number uint64
	Gen    uint64
}

func (v *Version) increment() {
	v.Number++
}

func (v1 *Version) replace(v2 Version) {
	v1.Number = v2.Number
}

func (v1 *Version) replaceAndIncrement(v2 Version) {
	v1.replace(v2)
	v1.increment()
}

func (v1 *Version) compare(v2 Version) int8 {
	if v1.Number > v2.Number {
		return 1
	}
	if v1.Number < v2.Number {
		return -1
	}
	return 0
}

func (n *Node) newVersion() {
	v := Version{1, 1}
	n.version = v
}

func (v *Version) touch() {
	v.increment()
}