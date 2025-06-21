package dogma

// DisableOption is an option that modifies the behavior of a disabled handler.
type DisableOption interface {
	futureDisableOption()
}
