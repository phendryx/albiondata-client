package photon_spectator

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

const fragmentBufferSize = 128

// Provides a LRU backed buffer which will assemble ReliableFragments
// into a single PhotonCommand with type ReliableMessage
type FragmentBuffer struct {
	cache *lru.Cache[int32, fragmentBufferEntry]
}

// Offers a message to the buffer. Returns nil when no new commands could be assembled from the
// buffer's contents.
func (buf *FragmentBuffer) Offer(msg ReliableFragment) *PhotonCommand {
	var entry fragmentBufferEntry

	if existing, ok := buf.cache.Get(msg.SequenceNumber); ok {
		entry = existing
		entry.Fragments[int(msg.FragmentNumber)] = msg.Data
	} else {
		entry.SequenceNumber = msg.SequenceNumber
		entry.FragmentsNeeded = int(msg.FragmentCount)
		entry.Fragments = make(map[int][]byte)
		entry.Fragments[int(msg.FragmentNumber)] = msg.Data
	}

	if entry.Finished() {
		command := entry.Make()
		buf.cache.Remove(msg.SequenceNumber)
		return &command
	} else {
		buf.cache.Add(msg.SequenceNumber, entry)
		return nil
	}
}

type fragmentBufferEntry struct {
	SequenceNumber  int32
	FragmentsNeeded int
	Fragments       map[int][]byte
}

func (buf fragmentBufferEntry) Finished() bool {
	return len(buf.Fragments) == buf.FragmentsNeeded
}

func (buf fragmentBufferEntry) Make() PhotonCommand {
	var data []byte

	for i := 0; i < buf.FragmentsNeeded; i++ {
		data = append(data, buf.Fragments[i]...)
	}

	return PhotonCommand{
		Type:                   SendReliableType,
		Data:                   data,
		ReliableSequenceNumber: buf.SequenceNumber,
	}
}

// Makes a new instance of a FragmentBuffer
func NewFragmentBuffer() *FragmentBuffer {
	var f FragmentBuffer
	// lru.New only returns an error if size <= 0, so this is safe
	f.cache, _ = lru.New[int32, fragmentBufferEntry](fragmentBufferSize)
	return &f
}
