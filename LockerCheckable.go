package golockercheckable

import (
	sync_o "sync"
)

type LockerCheckable interface {
	sync_o.Locker
	IsLocked() bool
	IsLocakedByMe() (locked bool, byme bool)
	LocekdByWho() (locked bool, goid uint64)
}
