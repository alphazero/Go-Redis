//   Copyright 2009-2012 Joubin Houshyar
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package redis

// Future types that wrap generic Redis response types - not entirely happy about
// this ...

// FutureKeys
//
type FutureKeys interface {
	Get() ([]string, Error)
	TryGet(timeout int64) (keys []string, error Error, ok bool)
}
type _futurekeys struct {
	future FutureBytes
}

func newFutureKeys(future FutureBytes) FutureKeys {
	return _futurekeys{future}
}
func (fvc _futurekeys) Get() (v []string, error Error) {
	gv, err := fvc.future.Get()
	if err != nil {
		return nil, err
	}
	v = convAndSplit(gv)
	return v, nil
}
func (fvc _futurekeys) TryGet(ns int64) (v []string, error Error, ok bool) {
	gv, err, ok := fvc.future.TryGet(ns)
	if !ok {
		return nil, nil, ok
	}
	if err != nil {
		return nil, err, ok
	}
	v = convAndSplit(gv)
	return v, nil, ok
}

// FutureInfo
//
type FutureInfo interface {
	Get() (map[string]string, Error)
	TryGet(timeout int64) (keys map[string]string, error Error, ok bool)
}
type _futureinfo struct {
	future FutureBytes
}

func newFutureInfo(future FutureBytes) FutureInfo {
	return _futureinfo{future}
}
func (fvc _futureinfo) Get() (v map[string]string, error Error) {
	gv, err := fvc.future.Get()
	if err != nil {
		return nil, err
	}
	v = parseInfo(gv)
	return v, nil
}
func (fvc _futureinfo) TryGet(ns int64) (v map[string]string, error Error, ok bool) {
	gv, err, ok := fvc.future.TryGet(ns)
	if !ok {
		return nil, nil, ok
	}
	if err != nil {
		return nil, err, ok
	}
	v = parseInfo(gv)
	return v, nil, ok
}

// FutureKeyType
//
type FutureKeyType interface {
	Get() (KeyType, Error)
	TryGet(timeout int64) (keys KeyType, error Error, ok bool)
}
type _futurekeytype struct {
	future FutureString
}

func newFutureKeyType(future FutureString) FutureKeyType {
	return _futurekeytype{future}
}
func (fvc _futurekeytype) Get() (v KeyType, error Error) {
	gv, err := fvc.future.Get()
	if err != nil {
		return RT_NONE, err
	}
	v = GetKeyType(gv)
	return v, nil
}
func (fvc _futurekeytype) TryGet(ns int64) (v KeyType, error Error, ok bool) {
	gv, err, ok := fvc.future.TryGet(ns)
	if !ok {
		return RT_NONE, nil, ok
	}
	if err != nil {
		return RT_NONE, err, ok
	}
	v = GetKeyType(gv)
	return v, nil, ok
}
