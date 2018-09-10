package maps

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
)

type Config struct {
	TagName  string
	OmitZero bool
	OmitNil  bool
}

var defaultConfig = &Config{
	TagName:  "map",
	OmitZero: false,
	OmitNil:  false,
}

// The below code is a lightly editied version of code written by the Go Authors
// for the encoding/json package. As such, it remains under the BSD-style
// license it was originally copywritten under.
// https://golang.org/LICENSE

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type field struct {
	name      string
	nameBytes []byte                 // []byte(name)
	equalFold func(s, t []byte) bool // bytes.EqualFold or equivalent

	tagged bool
	index  []int
	typ    reflect.Type

	omitZero bool
	omitNil  bool
	asValue  bool
}

func fillField(f field) field {
	f.nameBytes = []byte(f.name)
	f.equalFold = foldFunc(f.nameBytes)
	return f
}

// byIndex sorts field by index sequence.
type byIndex []field

func (fields byIndex) Len() int {
	return len(fields)
}

func (fields byIndex) Swap(i, j int) {
	fields[i], fields[j] = fields[j], fields[i]
}

func (fields byIndex) Less(i, j int) bool {
	l := fields[i]
	r := fields[j]
	for k := 0; k < len(l.index); k++ {
		if k >= len(r.index) {
			return false
		}
		if l.index[k] != r.index[k] {
			return l.index[k] < r.index[k]
		}
	}
	return len(l.index) < len(r.index)
}

var fieldCache struct {
	value atomic.Value // map[reflect.Type][]field
	mu    sync.Mutex   // used only by writers
}

// cachedTypeFields caches the return of typeFields to avoid repeated work.
func cachedTypeFields(t reflect.Type, cfg *Config) []field {
	m, _ := fieldCache.value.Load().(map[reflect.Type][]field)
	f := m[t]
	if f != nil {
		return f
	}

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = typeFields(t, cfg)
	if f == nil {
		f = []field{}
	}

	fieldCache.mu.Lock()
	m, _ = fieldCache.value.Load().(map[reflect.Type][]field)
	newM := make(map[reflect.Type][]field, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[t] = f
	fieldCache.value.Store(newM)
	fieldCache.mu.Unlock()
	return f
}

// typeFields returns a list of fields that should be recognized for the given
// type. The algorithm is breadth-first search over the set of structs to
// include - the top struct and then any reachable anonymous structs.
func typeFields(t reflect.Type, cfg *Config) []field {
	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				isUnexported := sf.PkgPath != ""
				isEmbedded := sf.Anonymous
				if isEmbedded {
					t := sf.Type
					if t.Kind() == reflect.Ptr {
						t = t.Elem()
					}
					if isUnexported && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-structs.
						// Would these be the fields of unexported Complex
						// numbers? I'm not sure when we would hit this.
						continue
					}
					// Do not ignore embedded fields of unexported structs,
					// because they may have exported fields.
				} else if isUnexported {
					continue
				}

				tag := sf.Tag.Get(cfg.TagName)
				tagged := tag != ""
				name, opts := parseTag(tag)
				if name == "-" {
					continue
				}
				// TODO: Consider adding a check to ensure `name` is a valid key
				if name == "" {
					name = sf.Name
				}
				omitZero := cfg.OmitZero
				if opts.Contain("omitZero") {
					omitZero = true
				} else if opts.Contain("noOmitZero") {
					omitZero = false
				}
				omitNil := cfg.OmitNil
				if opts.Contain("omitNil") {
					omitNil = true
				} else if opts.Contain("noOmitNil") {
					omitNil = false
				}
				asValue := false
				if opts.Contain("value") {
					asValue = true
				}

				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				sft := sf.Type
				// If we're looking at an unnamed pointer type, follow the
				// pointer to the underlying type.
				if sft.Name() == "" && sft.Kind() == reflect.Ptr {
					sft = sft.Elem()
				}

				// Record the found field and index sequence ...
				if tagged || !isEmbedded || sft.Kind() != reflect.Struct {
					fields = append(fields, fillField(field{
						name:     name,
						tagged:   tagged,
						index:    index,
						typ:      sft,
						omitZero: omitZero,
						omitNil:  omitNil,
						asValue:  asValue,
					}))
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second, so
						// that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// ... or record a new anonymous struct to be explored in the
				// next round.
				nextCount[sft]++
				if nextCount[sft] == 1 {
					next = append(next, fillField(field{
						name:  sft.Name(),
						index: index,
						typ:   sft,
					}))
				}
			}
		}
	}

	// Sort field first by name, then breaking ties with index sequence length,
	// then breaking ties with "name came from map tag", then breaking ties with
	// index sequence.
	sort.Slice(fields, func(i, j int) bool {
		l := fields[i]
		r := fields[j]

		if l.name != r.name {
			return l.name < r.name
		}
		if len(l.index) != len(r.index) {
			return len(l.index) < len(r.index)
		}
		if l.tagged != r.tagged {
			return l.tagged
		}
		for k := 0; k < len(l.index); k++ {
			if k >= len(r.index) {
				return false
			}
			if l.index[k] != r.index[k] {
				return l.index[k] < r.index[k]
			}
		}
		return len(l.index) < len(r.index)
	})

	// Delete all fields that are hidden based on Go's rules for embedded
	// fields, modified for the presence of map tags.
	out := fields[:0]
	for i, count := 0, 0; i < len(fields); i += count {
		// Find the number of fields that share a name
		for count = 1; i+count < len(fields); count++ {
			if fields[i].name != fields[i+count].name {
				break
			}
		}
		// If there's only one field named `name`, our job is easy ...
		if count == 1 {
			out = append(out, fields[i])
			continue
		}
		// ... otherwise, find the single field that dominates the other
		// similarly named fields using Go's embedding rules, modified by the
		// presence of map tags. If there are multiple top-level fields -- which
		// is an error in Go -- we mirror the compile-time "ambiguous selector"
		// error as closely as possible at runtime. With a panic.
		contended := fields[i : i+count]
		taggedIndex := -1
		for j, f := range contended {
			// A shorter index length indicates a more dominant field. The
			// `contended` slice is sorted in increasing index-length order. We
			// can therefore drop longer (and less dominant) entries by simply
			// truncating the slice.
			if len(f.index) > len(contended[0].index) {
				contended = contended[:j]
				break
			}
			// Fields with map tags are given special precedence for the given
			// index-length level, so we need to keep track of their presence.
			if f.tagged {
				// If there are multiple tagged fields at the same index level,
				// we have a genuine conflict.
				if taggedIndex >= 0 {
					panic(fmt.Errorf("encmap: ambiguous tagged field name '%s' in %s", f.name, t.Name()))
				}
				taggedIndex = j
			}
		}
		if taggedIndex >= 0 {
			out = append(out, contended[taggedIndex])
			continue
		}
		// All remaining contended fields have the same length. If there's more
		// than one, we have a conflict.
		if len(contended) > 1 {
			panic(fmt.Errorf("encmap: ambiguous tagged field name '%s' in %s", contended[0].name, t.Name()))
		}
		out = append(out, contended[0])
	}

	fields = out
	sort.Sort(byIndex(fields))

	return fields
}
