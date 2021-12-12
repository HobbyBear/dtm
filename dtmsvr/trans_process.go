/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dtmsvr

import (
	"time"

	"github.com/yedf/dtm/dtmcli"
	"github.com/yedf/dtm/dtmcli/dtmimp"
)

// Process process global transaction once
func (t *TransGlobal) Process() map[string]interface{} {
	r := t.process()
	transactionMetrics(t, r["dtm_result"] == dtmcli.ResultSuccess)
	return r
}

func (t *TransGlobal) process() map[string]interface{} {
	if t.Options != "" {
		dtmimp.MustUnmarshalString(t.Options, &t.TransOptions)
	}

	if !t.WaitResult {
		go t.processInner()
		return dtmcli.MapSuccess
	}
	submitting := t.Status == dtmcli.StatusSubmitted
	err := t.processInner()
	if err != nil {
		return map[string]interface{}{"dtm_result": dtmcli.ResultFailure, "message": err.Error()}
	}
	if submitting && t.Status != dtmcli.StatusSucceed {
		return map[string]interface{}{"dtm_result": dtmcli.ResultFailure, "message": "trans failed by user"}
	}
	return dtmcli.MapSuccess
}

func (t *TransGlobal) processInner() (rerr error) {
	defer handlePanic(&rerr)
	defer func() {
		if rerr != nil {
			dtmimp.LogRedf("processInner got error: %s", rerr.Error())
		}
		if TransProcessedTestChan != nil {
			dtmimp.Logf("processed: %s", t.Gid)
			TransProcessedTestChan <- t.Gid
			dtmimp.Logf("notified: %s", t.Gid)
		}
	}()
	dtmimp.Logf("processing: %s status: %s", t.Gid, t.Status)
	branches := []TransBranch{}
	dbGet().Must().Where("gid=?", t.Gid).Order("id asc").Find(&branches)
	t.lastTouched = time.Now()
	rerr = t.getProcessor().ProcessOnce(branches)
	return
}

func (t *TransGlobal) saveNew() error {
	branches := t.getProcessor().GenBranches()
	t.setNextCron(cronReset)
	t.Options = dtmimp.MustMarshalString(t.TransOptions)
	if t.Options == "{}" {
		t.Options = ""
	}
	return getStore().SaveNewTrans(&t.TransGlobalStore, branches)
}
