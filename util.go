package dogma

// nocmp is an embeddable type that ensures its parent isn't comparable.
type nocmp [0]func()
