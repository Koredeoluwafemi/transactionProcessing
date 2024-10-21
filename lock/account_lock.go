package lock

import "sync"

var accountLocks = make(map[uint]*sync.Mutex)
var lockMu sync.Mutex

// GetAccountLock ... Get the lock for the account
func GetAccountLock(accountID uint) *sync.Mutex {
	lockMu.Lock()
	defer lockMu.Unlock()

	if _, exists := accountLocks[accountID]; !exists {
		accountLocks[accountID] = &sync.Mutex{}
	}
	return accountLocks[accountID]
}
