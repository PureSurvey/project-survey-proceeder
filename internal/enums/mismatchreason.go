package enums

type MismatchReason int

const (
	MRUnknown MismatchReason = iota
	MRInvalidTracker
	MROutdatedTracker
)
