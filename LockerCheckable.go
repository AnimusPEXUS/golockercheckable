package golockercheckable

import (
	sync_o "sync"
)

type LockerCheckable interface {
	sync_o.Locker
	IsLocked() bool
}
