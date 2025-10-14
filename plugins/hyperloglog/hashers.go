// Copyright 2025 samber.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.apache.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rohyperloglog

import (
	// bearer:disable go_gosec_blocklist_md5
	"crypto/md5"
	// bearer:disable go_gosec_blocklist_sha1
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"hash/maphash"
)

// Convenience variables for easy access
var (
	// StringHash provides hash functions for strings
	StringHash = StringHashers{}
	// BytesHash provides hash functions for byte slices
	BytesHash = BytesHashers{}
)

// StringHashers provides common hash functions for strings
type StringHashers struct{}

// FNV64a returns a hash function using FNV-1a 64-bit
func (StringHashers) FNV64a() func(string) uint64 {
	return func(input string) uint64 {
		h := fnv.New64a()
		h.Write([]byte(input))
		return h.Sum64()
	}
}

// FNV64 returns a hash function using FNV-1 64-bit
func (StringHashers) FNV64() func(string) uint64 {
	return func(input string) uint64 {
		h := fnv.New64()
		h.Write([]byte(input))
		return h.Sum64()
	}
}

// FNV32a returns a hash function using FNV-1a 32-bit
func (StringHashers) FNV32a() func(string) uint64 {
	return func(input string) uint64 {
		h := fnv.New32a()
		h.Write([]byte(input))
		return uint64(h.Sum32())
	}
}

// FNV32 returns a hash function using FNV-1 32-bit
func (StringHashers) FNV32() func(string) uint64 {
	return func(input string) uint64 {
		h := fnv.New32()
		h.Write([]byte(input))
		return uint64(h.Sum32())
	}
}

// SHA256 returns a hash function using SHA-256
func (StringHashers) SHA256() func(string) uint64 {
	return func(input string) uint64 {
		hash := sha256.Sum256([]byte(input))
		return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
			uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
	}
}

// SHA1 returns a hash function using SHA-1
func (StringHashers) SHA1() func(string) uint64 {
	return func(input string) uint64 {
		// bearer:disable go_lang_weak_hash_sha1, go_gosec_crypto_weak_crypto
		hash := sha1.Sum([]byte(input))
		return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
			uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
	}
}

// SHA512 returns a hash function using SHA-512
func (StringHashers) SHA512() func(string) uint64 {
	return func(input string) uint64 {
		hash := sha512.Sum512([]byte(input))
		return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
			uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
	}
}

// MD5 returns a hash function using MD5
func (StringHashers) MD5() func(string) uint64 {
	return func(input string) uint64 {
		// bearer:disable go_lang_weak_hash_md5, go_gosec_crypto_weak_crypto
		hash := md5.Sum([]byte(input))
		return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
			uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
	}
}

// MapHash returns a hash function using maphash
func (StringHashers) MapHash() func(string) uint64 {
	return func(input string) uint64 {
		var h maphash.Hash
		h.WriteString(input)
		return h.Sum64()
	}
}

// CRC32 returns a hash function using CRC-32
func (StringHashers) CRC32() func(string) uint64 {
	return func(input string) uint64 {
		return uint64(crc32.ChecksumIEEE([]byte(input)))
	}
}

// CRC64 returns a hash function using CRC-64
func (StringHashers) CRC64() func(string) uint64 {
	return func(input string) uint64 {
		return crc64.Checksum([]byte(input), crc64.MakeTable(crc64.ISO))
	}
}

// Adler32 returns a hash function using Adler-32
func (StringHashers) Adler32() func(string) uint64 {
	return func(input string) uint64 {
		return uint64(adler32.Checksum([]byte(input)))
	}
}

// Jenkins returns a hash function using Jenkins hash (one-at-a-time)
func (StringHashers) Jenkins() func(string) uint64 {
	return func(input string) uint64 {
		var hash uint64
		for _, c := range input {
			hash += uint64(c)
			hash += (hash << 10)
			hash ^= (hash >> 6)
		}
		hash += (hash << 3)
		hash ^= (hash >> 11)
		hash += (hash << 15)
		return hash
	}
}

// DJB2 returns a hash function using DJB2 algorithm
func (StringHashers) DJB2() func(string) uint64 {
	return func(input string) uint64 {
		var hash uint64 = 5381
		for _, c := range input {
			hash = ((hash << 5) + hash) + uint64(c) // hash * 33 + c
		}
		return hash
	}
}

// SDBM returns a hash function using SDBM algorithm
func (StringHashers) SDBM() func(string) uint64 {
	return func(input string) uint64 {
		var hash uint64 = 0
		for _, c := range input {
			hash = uint64(c) + (hash << 6) + (hash << 16) - hash
		}
		return hash
	}
}

// Loselose returns a hash function using the "lose lose" algorithm
func (StringHashers) Loselose() func(string) uint64 {
	return func(input string) uint64 {
		var hash uint64 = 0
		for _, c := range input {
			hash += uint64(c)
		}
		return hash
	}
}

// BytesHashers provides common hash functions for byte slices
type BytesHashers struct{}

// FNV64a returns a hash function using FNV-1a 64-bit for byte slices
func (BytesHashers) FNV64a() func([]byte) uint64 {
	return func(input []byte) uint64 {
		h := fnv.New64a()
		h.Write(input)
		return h.Sum64()
	}
}

// FNV64 returns a hash function using FNV-1 64-bit for byte slices
func (BytesHashers) FNV64() func([]byte) uint64 {
	return func(input []byte) uint64 {
		h := fnv.New64()
		h.Write(input)
		return h.Sum64()
	}
}

// FNV32a returns a hash function using FNV-1a 32-bit for byte slices
func (BytesHashers) FNV32a() func([]byte) uint64 {
	return func(input []byte) uint64 {
		h := fnv.New32a()
		h.Write(input)
		return uint64(h.Sum32())
	}
}

// FNV32 returns a hash function using FNV-1 32-bit for byte slices
func (BytesHashers) FNV32() func([]byte) uint64 {
	return func(input []byte) uint64 {
		h := fnv.New32()
		h.Write(input)
		return uint64(h.Sum32())
	}
}

// SHA256 returns a hash function using SHA-256 for byte slices
func (BytesHashers) SHA256() func([]byte) uint64 {
	return func(input []byte) uint64 {
		hash := sha256.Sum256(input)
		return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
			uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
	}
}

// SHA1 returns a hash function using SHA-1 for byte slices
func (BytesHashers) SHA1() func([]byte) uint64 {
	return func(input []byte) uint64 {
		// bearer:disable go_lang_weak_hash_sha1, go_gosec_crypto_weak_crypto
		hash := sha1.Sum(input)
		return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
			uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
	}
}

// SHA512 returns a hash function using SHA-512 for byte slices
func (BytesHashers) SHA512() func([]byte) uint64 {
	return func(input []byte) uint64 {
		hash := sha512.Sum512(input)
		return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
			uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
	}
}

// MD5 returns a hash function using MD5 for byte slices
func (BytesHashers) MD5() func([]byte) uint64 {
	return func(input []byte) uint64 {
		// bearer:disable go_lang_weak_hash_md5, go_gosec_crypto_weak_crypto
		hash := md5.Sum(input)
		return uint64(hash[0])<<56 | uint64(hash[1])<<48 | uint64(hash[2])<<40 | uint64(hash[3])<<32 |
			uint64(hash[4])<<24 | uint64(hash[5])<<16 | uint64(hash[6])<<8 | uint64(hash[7])
	}
}

// MapHash returns a hash function using maphash for byte slices
func (BytesHashers) MapHash() func([]byte) uint64 {
	return func(input []byte) uint64 {
		var h maphash.Hash
		h.Write(input)
		return h.Sum64()
	}
}

// CRC32 returns a hash function using CRC-32 for byte slices
func (BytesHashers) CRC32() func([]byte) uint64 {
	return func(input []byte) uint64 {
		return uint64(crc32.ChecksumIEEE(input))
	}
}

// CRC64 returns a hash function using CRC-64 for byte slices
func (BytesHashers) CRC64() func([]byte) uint64 {
	return func(input []byte) uint64 {
		return crc64.Checksum(input, crc64.MakeTable(crc64.ISO))
	}
}

// Adler32 returns a hash function using Adler-32 for byte slices
func (BytesHashers) Adler32() func([]byte) uint64 {
	return func(input []byte) uint64 {
		return uint64(adler32.Checksum(input))
	}
}

// Jenkins returns a hash function using Jenkins hash (one-at-a-time) for byte slices
func (BytesHashers) Jenkins() func([]byte) uint64 {
	return func(input []byte) uint64 {
		var hash uint64
		for _, c := range input {
			hash += uint64(c)
			hash += (hash << 10)
			hash ^= (hash >> 6)
		}
		hash += (hash << 3)
		hash ^= (hash >> 11)
		hash += (hash << 15)
		return hash
	}
}

// DJB2 returns a hash function using DJB2 algorithm for byte slices
func (BytesHashers) DJB2() func([]byte) uint64 {
	return func(input []byte) uint64 {
		var hash uint64 = 5381
		for _, c := range input {
			hash = ((hash << 5) + hash) + uint64(c) // hash * 33 + c
		}
		return hash
	}
}

// SDBM returns a hash function using SDBM algorithm for byte slices
func (BytesHashers) SDBM() func([]byte) uint64 {
	return func(input []byte) uint64 {
		var hash uint64 = 0
		for _, c := range input {
			hash = uint64(c) + (hash << 6) + (hash << 16) - hash
		}
		return hash
	}
}

// Loselose returns a hash function using the "lose lose" algorithm for byte slices
func (BytesHashers) Loselose() func([]byte) uint64 {
	return func(input []byte) uint64 {
		var hash uint64 = 0
		for _, c := range input {
			hash += uint64(c)
		}
		return hash
	}
}
