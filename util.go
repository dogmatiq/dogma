package dogma

// nocmp is an embeddable type that ensures its parent is not comparable.
type nocmp [0]func()
