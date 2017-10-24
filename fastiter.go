package roaring

type FastIterator struct {
	pos              int
	arraypos         int
	carray           []uint32
	highlowcontainer *roaringArray
}

func (ii *FastIterator) HasNext() bool {
	return ii.pos < ii.highlowcontainer.size()
}

func (ii *FastIterator) Next() uint32 {
	if ii.arraypos >= len(ii.carray) {
		ii.arraypos = 0

		hs := uint32(ii.highlowcontainer.getKeyAtIndex(ii.pos)) << 16
		c := ii.highlowcontainer.getContainerAtIndex(ii.pos)
		ii.pos++

		csize := c.getCardinality()
		if len(ii.carray) < csize {
			ii.carray = make([]uint32, csize)
		} else {
			ii.carray = ii.carray[:csize]
		}

		c.fillLeastSignificant16bits(ii.carray, 0, hs)
	}

	ret := ii.carray[ii.arraypos]
	ii.arraypos++
	return ret
}

func newFastIterator(a *Bitmap) *FastIterator {
	p := new(FastIterator)
	p.arraypos = 0
	p.pos = 0
	p.highlowcontainer = &a.highlowcontainer
	return p
}

// Iterator creates a new IntIterable to iterate over the integers contained in the bitmap, in sorted order
func (rb *Bitmap) FastIterator() *FastIterator {
	return newFastIterator(rb)
}

// FillArray fills a slice containing all of the integers stored in the Bitmap in sorted order
func (rb *Bitmap) FillArray(array []uint32) {
	pos := 0
	pos2 := 0

	for pos < rb.highlowcontainer.size() {
		hs := uint32(rb.highlowcontainer.getKeyAtIndex(pos)) << 16
		c := rb.highlowcontainer.getContainerAtIndex(pos)
		pos++
		c.fillLeastSignificant16bits(array, pos2, hs)
		pos2 += c.getCardinality()
	}
}
