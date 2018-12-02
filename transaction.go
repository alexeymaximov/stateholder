package stateholder

// Begin transaction.
func (sh *Stateholder) Begin() error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping == nil {
		return &ErrorDetached{}
	}
	if sh.transaction {
		return &ErrorTransactionAlreadyStarted{}
	}
	sh.transaction = true
	return nil
}

// Rollback transaction.
func (sh *Stateholder) Rollback() error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping == nil {
		return &ErrorDetached{}
	}
	if !sh.transaction {
		return &ErrorTransactionNotStarted{}
	}
	for _, entry := range sh.entries {
		entry.buffer = nil
	}
	sh.transaction = false
	return nil
}

// Commit transaction.
func (sh *Stateholder) commit() error {
	if !sh.transaction {
		return &ErrorTransactionNotStarted{}
	}
	// TODO: Add full commit buffer.
	for _, entry := range sh.entries {
		if entry.buffer != nil {
			if n, err := sh.mapping.WriteAt(entry.buffer, entry.offset); err != nil {
				return err
			} else if n != int(entry.size) {
				return &ErrorCorruptedWrite{Real: n, Expected: int(entry.size)}
			}
			entry.buffer = nil
		}
	}
	sh.transaction = false
	return nil
}

// Commit transaction.
func (sh *Stateholder) Commit() error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping == nil {
		return &ErrorDetached{}
	}
	return sh.commit()
}

// Commit transaction and sync data.
func (sh *Stateholder) Persist() error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping == nil {
		return &ErrorDetached{}
	}
	if err := sh.commit(); err != nil {
		return err
	}
	return sh.mapping.Sync()
}
