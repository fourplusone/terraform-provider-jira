// Copyright 2015-2018 trivago N.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package treflect

import (
	"testing"

	"github.com/trivago/tgo/ttesting"
)

type CloneTest1 struct {
	Number  int
	Text    string
	Slice   []int
	Array   [3]int
	private int // private fields should be ignored
}

type cloneTest2 struct {
	CloneTest1
	Slice2 []CloneTest1
	Dict   map[string]CloneTest1
}

func TestClone(t *testing.T) {
	expect := ttesting.NewExpect(t)

	test1 := CloneTest1{
		Number: 1,
		Text:   "test",
		Slice:  []int{1, 2, 3},
		Array:  [3]int{1, 2, 3},
	}

	testClone1 := Clone(test1)
	expect.Equal(test1, testClone1)

	test2 := cloneTest2{
		CloneTest1: test1,
		Slice2:     []CloneTest1{test1},
		Dict: map[string]CloneTest1{
			"foo": test1,
		},
	}

	testClone2 := Clone(test2)
	expect.Equal(test2, testClone2)
}
