/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package test

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yedf/dtm/common"
	"github.com/yedf/dtm/dtmcli"
	"github.com/yedf/dtm/dtmcli/dtmimp"
	"github.com/yedf/dtm/dtmsvr"
	"github.com/yedf/dtm/dtmsvr/storage"
	"github.com/yedf/dtm/examples"
)

var DtmServer = examples.DtmHttpServer
var Busi = examples.Busi
var app *gin.Engine

func getTransStatus(gid string) string {
	sm := TransGlobal{}
	err := storage.GetStore().GetTransGlobal(gid, &sm.TransGlobalStore)
	e2p(err)
	return sm.Status
}

func getBranchesStatus(gid string) []string {
	branches := storage.GetStore().GetBranches(gid)
	status := []string{}
	for _, branch := range branches {
		status = append(status, branch.Status)
	}
	return status
}

func assertSucceed(t *testing.T, gid string) {
	waitTransProcessed(gid)
	assert.Equal(t, StatusSucceed, getTransStatus(gid))
}

func TestUpdateBranchAsync(t *testing.T) {
	common.DtmConfig.UpdateBranchSync = 0
	saga := genSaga1(dtmimp.GetFuncName(), false, false)
	saga.SetOptions(&dtmcli.TransOptions{WaitResult: true})
	err := saga.Submit()
	assert.Nil(t, err)
	waitTransProcessed(saga.Gid)
	time.Sleep(dtmsvr.UpdateBranchAsyncInterval)
	assert.Equal(t, []string{StatusPrepared, StatusSucceed}, getBranchesStatus(saga.Gid))
	assert.Equal(t, StatusSucceed, getTransStatus(saga.Gid))
	common.DtmConfig.UpdateBranchSync = 1
}
