// Copyright (c) 2012, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package storage provides storage abstraction for LevelDB.
package storage

import (
	"errors"
	"io"
)

type FileType uint32

const (
	TypeManifest FileType = 1 << iota
	TypeJournal
	TypeTable

	TypeAll = TypeManifest | TypeJournal | TypeTable
)

func (t FileType) String() string {
	switch t {
	case TypeManifest:
		return "manifest"
	case TypeJournal:
		return "journal"
	case TypeTable:
		return "table"
	}
	return "<unknown>"
}

var (
	ErrInvalidFile = errors.New("invalid file for argument")
	ErrLocked      = errors.New("already locked")
	ErrNotLocked   = errors.New("not locked")
	ErrInvalidLock = errors.New("invalid lock handle")
)

type Syncer interface {
	Sync() error
}

type Reader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type Writer interface {
	io.Writer
	io.Closer
	Syncer
}

type Locker interface {
	Release() error
}

type File interface {
	// Open file for read.
	// Should return os.ErrNotExist if the file does not exist.
	Open() (r Reader, err error)

	// Create file for write. Truncate if file already exist.
	Create() (w Writer, err error)

	// Rename to given number and type.
	Rename(number uint64, t FileType) error

	// Return true if the file is exist.
	Exist() bool

	// Return file type.
	Type() FileType

	// Return file number
	Num() uint64

	// Return size of the file.
	Size() (size uint64, err error)

	// Remove file.
	Remove() error
}

type Storage interface {
	// Lock the storage, so any subsequent attempt to lock the same storage
	// will fail.
	Lock() (l Locker, err error)

	// Print a string, for logging.
	Print(str string)

	// Get file with given number and type.
	GetFile(number uint64, t FileType) File

	// Get all files that match given file types; multiple file type
	// may OR'ed together.
	GetFiles(t FileType) []File

	// Get manifest file.
	// Should return os.ErrNotExist if there's no current manifest file.
	GetManifest() (f File, err error)

	// Set manifest to given file.
	SetManifest(f File) error
}
