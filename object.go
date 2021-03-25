package log

/*
   Copyright 2019 Bruno Moura <brunotm@gmail.com>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

var (
	nullBytes = []byte(`null`)
)

// Object value
type Object struct {
	enc *encoder
}

// Bool adds the given bool key/value
func (o Object) Bool(key string, value bool) (object Object) {
	o.enc.addKey(key)
	o.enc.AppendBool(value)
	return o
}

// Float64 adds the given float key/value
func (o Object) Float64(key string, value float64) (object Object) {
	o.enc.addKey(key)
	o.enc.AppendFloat64(value)
	return o
}

// Int64 adds the given int key/value
func (o Object) Int64(key string, value int64) (object Object) {
	o.enc.addKey(key)
	o.enc.AppendInt64(value)
	return o
}

// Uint64 adds the given uint key/value
func (o Object) Uint64(key string, value uint64) (object Object) {
	o.enc.addKey(key)
	o.enc.AppendUint64(value)
	return o
}

// String adds the given string key/value
func (o Object) String(key string, value string) (object Object) {
	o.enc.addKey(key)
	o.enc.AppendString(value)
	return o
}

// Null adds a null value for the given key
func (o Object) Null(key string) (object Object) {
	o.enc.addKey(key)
	o.enc.AppendBytes(nullBytes)
	return o
}

// Error adds a error value for the given key
func (o Object) Error(key string, err error) (object Object) {
	if err == nil {
		return o.Null(key)
	}

	return o.String(key, err.Error())
}
