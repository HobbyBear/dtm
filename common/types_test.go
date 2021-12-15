/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yedf/dtm/dtmcli/dtmimp"
)

func TestGeneralDB(t *testing.T) {
	MustLoadConfig()
	if DtmConfig.DB["driver"] == "redis" {

	} else {
		testSql(t)
		testWaitDBUp(t)
		testDbAlone(t)
	}
}
func testSql(t *testing.T) {
	db := DbGet(DtmConfig.DB)
	err := func() (rerr error) {
		defer dtmimp.P2E(&rerr)
		dbr := db.NoMust().Exec("select a")
		assert.NotEqual(t, nil, dbr.Error)
		db.Must().Exec("select a")
		return nil
	}()
	assert.NotEqual(t, nil, err)
}

func testWaitDBUp(t *testing.T) {
	WaitDBUp()
}

func testDbAlone(t *testing.T) {
	db, err := dtmimp.StandaloneDB(DtmConfig.DB)
	assert.Nil(t, err)
	_, err = dtmimp.DBExec(db, "select 1")
	assert.Equal(t, nil, err)
	_, err = dtmimp.DBExec(db, "")
	assert.Equal(t, nil, err)
	db.Close()
	_, err = dtmimp.DBExec(db, "select 1")
	assert.NotEqual(t, nil, err)
}

func TestConfig(t *testing.T) {
	testConfigStringField(DtmConfig.DB, "driver", "", t)
	testConfigStringField(DtmConfig.DB, "user", "", t)
	testConfigIntField(&DtmConfig.RetryInterval, 9, t)
	testConfigIntField(&DtmConfig.TimeoutToFail, 9, t)
}

func testConfigStringField(m map[string]string, key string, val string, t *testing.T) {
	old := m[key]
	m[key] = val
	str := checkConfig()
	assert.NotEqual(t, "", str)
	m[key] = old
}

func testConfigIntField(fd *int64, val int64, t *testing.T) {
	old := *fd
	*fd = val
	str := checkConfig()
	assert.NotEqual(t, "", str)
	*fd = old
}
