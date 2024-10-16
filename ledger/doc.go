// Package ledger provides a convenience cache over the Store, providing access by ID to Todo objects.
// Uses the Store as backend to ensure durability of the data blobs.
// Ledger is meant to be higher level than the durable store: the store deals with data blobs,
// while the Ledger deals with Objects, serializing/deserializing them when communicating with the store.
package ledger
