/*
 * Minio Client (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/minio/mc/internal/github.com/minio/minio/pkg/probe"
	. "github.com/minio/mc/internal/gopkg.in/check.v1"
	"github.com/minio/mc/pkg/console"
)

func (s *CmdTestSuite) TestDiffObjects(c *C) {
	/// filesystem
	root1, err := ioutil.TempDir(os.TempDir(), "cmd-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root1)

	root2, err := ioutil.TempDir(os.TempDir(), "cmd-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root2)

	objectPath1 := filepath.Join(root1, "object1")
	data := "hello"
	dataLen := len(data)
	perr := putTarget(objectPath1, int64(dataLen), bytes.NewReader([]byte(data)))
	c.Assert(perr, IsNil)

	objectPath2 := filepath.Join(root2, "object1")
	data = "hello"
	dataLen = len(data)
	perr = putTarget(objectPath2, int64(dataLen), bytes.NewReader([]byte(data)))
	c.Assert(perr, IsNil)

	for diff := range doDiffCmd(objectPath1, objectPath2, false) {
		c.Assert(diff.err, IsNil)
		c.Assert(len(diff.message), Equals, 0)
	}
}

func (s *CmdTestSuite) TestDiffDirs(c *C) {
	/// filesystem
	root1, err := ioutil.TempDir(os.TempDir(), "cmd-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root1)

	root2, err := ioutil.TempDir(os.TempDir(), "cmd-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root2)

	var perr *probe.Error
	for i := 0; i < 10; i++ {
		objectPath := filepath.Join(root1, "object"+strconv.Itoa(i))
		data := "hello"
		dataLen := len(data)
		perr = putTarget(objectPath, int64(dataLen), bytes.NewReader([]byte(data)))
		c.Assert(perr, IsNil)
	}

	for i := 0; i < 10; i++ {
		objectPath := filepath.Join(root2, "object"+strconv.Itoa(i))
		data := "hello"
		dataLen := len(data)
		perr = putTarget(objectPath, int64(dataLen), bytes.NewReader([]byte(data)))
		c.Assert(perr, IsNil)
	}

	for diff := range doDiffCmd(root1, root2, false) {
		c.Assert(diff.err, IsNil)
		c.Assert(len(diff.message), Equals, 0)
	}
}

func (s *CmdTestSuite) TestDiffContext(c *C) {
	err := app.Run([]string{os.Args[0], "diff", server.URL + "/bucket", server.URL + "/bucket"})
	c.Assert(err, IsNil)
	c.Assert(console.IsExited, Equals, false)

	err = app.Run([]string{os.Args[0], "diff", server.URL + "/bucket...", server.URL + "/bucket"})
	c.Assert(err, IsNil)
	c.Assert(console.IsExited, Equals, false)

	err = app.Run([]string{os.Args[0], "diff", server.URL + "/invalid", server.URL + "/invalid..."})
	c.Assert(err, IsNil)
	c.Assert(console.IsExited, Equals, true)
	// reset back
	console.IsExited = false
}
