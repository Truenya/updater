package util

type UpdaterType uint8

const (
	UpdateLocal UpdaterType = iota
	UpdateRemote
	UpdateInPlace
)
