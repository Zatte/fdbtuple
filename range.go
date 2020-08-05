package fdbtuple

// A Range describes all keys between a begin (inclusive) and end (exclusive)
// key selector.
type Range interface {
	// FDBRangeKeySelectors returns a pair of key selectors that describe the
	// beginning and end of a range.
	FDBRangeKeySelectors() (begin, end Selectable)
}

// An ExactRange describes all keys between a begin (inclusive) and end
// (exclusive) key. If you need to specify an ExactRange and you have only a
// Range, you must resolve the selectors returned by
// (Range).FDBRangeKeySelectors to keys using the (Transaction).GetKey method.
//
// Any object that implements ExactRange also implements Range, and may be used
// accordingly.
type ExactRange interface {
	// FDBRangeKeys returns a pair of keys that describe the beginning and end
	// of a range.
	FDBRangeKeys() (begin, end KeyConvertible)

	// An object that implements ExactRange must also implement Range
	// (logically, by returning FirstGreaterOrEqual of the keys returned by
	// FDBRangeKeys).
	Range
}
