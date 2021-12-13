/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dtmsvr

import (
	"time"

	"github.com/google/uuid"
	"github.com/yedf/dtm/common"
	"github.com/yedf/dtm/dtmcli/dtmimp"
	"github.com/yedf/dtm/dtmsvr/storage"
)

type branchStatus struct {
	id         uint64
	status     string
	finishTime *time.Time
}

var p2e = dtmimp.P2E
var e2p = dtmimp.E2P

var config = &common.DtmConfig

func getStore() storage.Store {
	return storage.GetStore()
}

// TransProcessedTestChan only for test usage. when transaction processed once, write gid to this chan
var TransProcessedTestChan chan string = nil

// GenGid generate gid, use uuid
func GenGid() string {
	return uuid.NewString()
}

// transFromDb construct trans from db
func transFromDb(gid string) *TransGlobal {
	m := TransGlobal{}
	err := getStore().GetTransGlobal(gid, &m.TransGlobalStore)
	if err == storage.ErrNotFound {
		return nil
	}
	e2p(err)
	return &m
}
