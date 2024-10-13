package store

// Package store implements a minimal persistent  key-value store whose
// keys are auto-increment keys and whose values are data blobs, aka slice of bytes.
// Implements backend to
// - redis
// - filesystem directory
