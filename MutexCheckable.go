package golockercheckable

import (
	sync_o "sync"

	"github.com/AnimusPEXUS/goroutineid"
)

var _ LockerCheckable = &MutexCheckable{}

// Difference to sync.Mutex, is what Unlock can be called if already unlocked.
// This can be used in sync.Cond without afraiding to get error when Unlock()
// called.
// NewMutexCheckable() must be used to create object of this struct.
type MutexCheckable struct {
	mutex_o   *sync_o.Mutex
	is_locked bool
	goid      uint64
	s         *sync_o.Mutex
}

// set locked to true, to make resulting Mutex already locked on creation
func NewMutexCheckable(locked bool) *MutexCheckable {
	self := &MutexCheckable{
		is_locked: locked,
		mutex_o:   &sync_o.Mutex{},
		s:         &sync_o.Mutex{},
	}
	if locked {
		self.mutex_o.Lock()
	}
	return self
}

// same as sync.Lock()
func (self *MutexCheckable) Lock() {

	// to call or not to call self.s.Unlock() in defer
	var already_unlocked bool = false

	self.s.Lock()
	defer func() {
		if !already_unlocked {
			self.s.Unlock()
		}
	}()

	// to make object functions available
	if self.is_locked {
		// self.s.Unlock() have to be called to allow calls to other functions
		// of this object.
		// this block available only if MutexCheckable already locked, so
		// this is ok to call self.s.Unlock() now, as self.mutex_o.Lock()
		// will already provide race safety
		already_unlocked = true
		self.s.Unlock()
	}

	self.is_locked = true
	id, err := goroutineid.GetCurrentGoId_byRuntimeStack()
	if err != nil {
		panic("can't get goroutine id")
	}
	self.goid = id
	self.mutex_o.Lock()
}

// same as sync.Unlock(), except, doesn't results in error when called on
// unlocked MutexCheckable
func (self *MutexCheckable) Unlock() {

	self.s.Lock()
	defer self.s.Unlock()

	if !self.is_locked {
		return
	}

	self.mutex_o.Unlock()

	self.is_locked = false
}

// returns lock state
func (self *MutexCheckable) IsLocked() bool {
	self.s.Lock()
	defer self.s.Unlock()
	return self.is_locked
}

func (self *MutexCheckable) IsLocakedByMe() (locked bool, byme bool) {
	self.s.Lock()
	defer self.s.Unlock()

	id, err := goroutineid.GetCurrentGoId_byRuntimeStack()
	if err != nil {
		panic("can't get goroutine id")
	}

	if self.is_locked {
		return true, id == self.goid
	} else {
		return false, false
	}
}

func (self *MutexCheckable) LocekdByWho() (locked bool, goid uint64) {
	self.s.Lock()
	defer self.s.Unlock()

	if self.is_locked {
		return true, self.goid
	} else {
		return false, 0
	}
}
