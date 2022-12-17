package storage

import (
	"testing"
)

func CommonStorageTest(storage Storage, t *testing.T) {
	if storage == nil {
		// nil without any test,
		return
	}

	// clean value
	err := storage.Clean()
	if err != nil {
		t.Error(err)
		return
	}

	if actSize, actLen, isClean := isStorageClean(storage); !isClean {
		t.Errorf("Storage already call the clean, but it not work or Size() or Len() not work: "+
			"actual size: %d; actual len: %d", actSize, actLen)
		return
	}

	knownKey := "knownKey"
	knownValue := "knownValue"
	samePrefix := "same_prefix"
	changedValue := "changed"

	baseValue := map[string]string{
		samePrefix + "/1": "1-value",
		samePrefix + "/2": "2-value",
		knownKey:          knownValue,
	}

	currentSize := uint64(0)
	currentLen := uint64(0)

	for key, value := range baseValue {
		storage.Set(key, value)

		currentSize += uint64(len([]byte(key))) + uint64(len([]byte(value)))
		currentLen++
	}

	// check the len and size
	if actLen, ok := isStorageLenExpected(storage, currentLen); !ok {
		t.Errorf("After set value, expect len: %d, but got is: %d", currentLen, actLen)
		return
	}

	if actSize, ok := isStorageSizeExpected(storage, currentSize); !ok {
		t.Errorf("After set value, expect size: %d, but got is: %d", len(baseValue), actSize)
		return
	}

	// check all path is set, and value is right
	for key, value := range baseValue {
		actualValue, ok := storage.Get(key)
		if !ok {
			t.Errorf("The path: [%s] is set by Set(), but can't get.", key)
			return
		}
		if actualValue != value {
			t.Errorf("The path: [%s] set the value: [%s], but got: [%s]", key, value, actualValue)
			return
		}
	}

	if ok := storage.Put(knownKey, changedValue); ok {
		t.Errorf("Specify path [%s] is parsent, but put success", knownKey)
		return
	}

	storage.Set(knownKey, changedValue)
	currentSize = currentSize - uint64(len([]byte(knownValue))) + uint64(len([]byte(changedValue)))

	if actValue, ok := storage.Get(knownKey); !ok {
		if !ok {
			t.Errorf("The path: [%s] is set by Set(), but can't get.", knownKey)
			return
		}
		if actValue != changedValue {
			t.Errorf("The path: [%s] set the value: [%s], but got: [%s]", knownKey, changedValue, actValue)
			return
		}
	}

	// check the len and size
	if actLen, ok := isStorageLenExpected(storage, currentLen); !ok {
		t.Errorf("After set value, expect len: %d, but got is: %d", currentLen, actLen)
		return
	}

	if actSize, ok := isStorageSizeExpected(storage, currentSize); !ok {
		t.Errorf("After set value, expect size: %d, but got is: %d", len(baseValue), actSize)
		return
	}
	// TODO delete clean contain keys
}

func isStorageClean(s Storage) (size, len uint64, isCleaned bool) {
	actSize, sizeIsZero := isStorageSizeExpected(s, 0)
	actLen, lenIsZero := isStorageLenExpected(s, 0)
	return actSize, actLen, sizeIsZero && lenIsZero
}

func isStorageLenExpected(s Storage, lenV uint64) (actual uint64, ok bool) {
	actual = s.Len()
	ok = actual == lenV
	return
}

func isStorageSizeExpected(s Storage, sizeV uint64) (actual uint64, ok bool) {
	actual = s.Size()
	ok = actual == sizeV
	return
}
