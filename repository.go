package main

type Repository interface {
	LongOperation() error
}
