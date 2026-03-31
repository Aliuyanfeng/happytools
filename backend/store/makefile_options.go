// Package store
// @description: Makefile Visual Editor — bbolt CRUD helpers for recent files and custom templates
package store

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

const (
	makefileRecentKey = "list"
	maxRecentFiles    = 10
)

// SaveRecentFile inserts path at the head of the recent-files list,
// deduplicates, and truncates to maxRecentFiles entries.
func SaveRecentFile(path string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(makefileRecentBucket)

		// Load existing list
		existing := make([]string, 0)
		if v := b.Get([]byte(makefileRecentKey)); v != nil {
			if err := json.Unmarshal(v, &existing); err != nil {
				return err
			}
		}

		// Deduplicate: remove any existing occurrence of path
		deduped := make([]string, 0, len(existing))
		for _, p := range existing {
			if p != path {
				deduped = append(deduped, p)
			}
		}

		// Prepend new path
		updated := append([]string{path}, deduped...)

		// Truncate to max
		if len(updated) > maxRecentFiles {
			updated = updated[:maxRecentFiles]
		}

		data, err := json.Marshal(updated)
		if err != nil {
			return err
		}
		return b.Put([]byte(makefileRecentKey), data)
	})
}

// GetRecentFiles returns the list of recently opened Makefile paths.
func GetRecentFiles() ([]string, error) {
	paths := make([]string, 0)
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(makefileRecentBucket)
		v := b.Get([]byte(makefileRecentKey))
		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &paths)
	})
	return paths, err
}

// RemoveRecentFile removes a specific path from the recent-files list.
func RemoveRecentFile(path string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(makefileRecentBucket)
		existing := make([]string, 0)
		if v := b.Get([]byte(makefileRecentKey)); v != nil {
			if err := json.Unmarshal(v, &existing); err != nil {
				return err
			}
		}
		filtered := make([]string, 0, len(existing))
		for _, p := range existing {
			if p != path {
				filtered = append(filtered, p)
			}
		}
		data, err := json.Marshal(filtered)
		if err != nil {
			return err
		}
		return b.Put([]byte(makefileRecentKey), data)
	})
}

// SaveMakefileTemplate stores a custom template (as raw JSON) by its ID.
// The caller is responsible for marshaling the template to JSON.
func SaveMakefileTemplate(id string, data []byte) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(makefileTemplateBucket)
		return b.Put([]byte(id), data)
	})
}

// GetMakefileTemplates returns all custom templates as raw JSON blobs keyed by ID.
func GetMakefileTemplates() ([][]byte, error) {
	var blobs [][]byte
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(makefileTemplateBucket)
		return b.ForEach(func(k, v []byte) error {
			cp := make([]byte, len(v))
			copy(cp, v)
			blobs = append(blobs, cp)
			return nil
		})
	})
	return blobs, err
}

// DeleteCustomTemplate removes a custom template by ID.
func DeleteCustomTemplate(id string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(makefileTemplateBucket)
		return b.Delete([]byte(id))
	})
}
