// Package dogma is a specification and API for message-based applications.
//
// It attempts to define a practical standard for authoring message-based
// applications in a manner agnostic to the mechanisms by which messages are
// transported and application state is persisted.
//
// To that end, this package defines a set of interfaces and recommendations for
// both "application developers" - those who are authoring a message-based
// application, and for "engine developers" - those who are authoring the
// platform responsible for message transport and persistence.
//
// The API documentation takes the form of a specification, and as such the key
// words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD
// NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
// interpreted as described in RFC 2119.
package dogma
