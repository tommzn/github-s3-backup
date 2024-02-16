package main

type Command interface {
	WithOptions(options Options)

	Exec() error

	Verify() error
}
