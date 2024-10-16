// Package model implements (a close-enough variant of) the Model part of the Model-View-Controller (MVC) pattern.
// model object(s), currently only the Todo object, implement the data itself and add behavior to represent
// the valid transformation of the data.
// The state transitions are implemented in this layer, close as possible to the data.
// Higher layers, like the controllers, will manipulate the objects based on the allowed operations.
package model
