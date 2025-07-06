package dogma

// DisableOption is an option that affects the behavior of a disabled handler.
type DisableOption interface {
	futureDisableOption()
}
