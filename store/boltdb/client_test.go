/**
 * Tencent is pleased to support the open source community by making Polaris available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package boltdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"

	api "github.com/polarismesh/polaris-server/common/api/v1"
	"github.com/polarismesh/polaris-server/common/model"
)

func Test_ConvertToClientObject(t *testing.T) {
	client := &api.Client{
		Host:    &wrapperspb.StringValue{Value: "1"},
		Type:    0,
		Version: &wrapperspb.StringValue{Value: "1"},
		Location: &api.Location{
			Region: &wrapperspb.StringValue{Value: "1"},
			Zone:   &wrapperspb.StringValue{Value: "1"},
			Campus: &wrapperspb.StringValue{Value: "1"},
		},
		Id: &wrapperspb.StringValue{Value: "1"},
		Stat: []*api.StatInfo{
			{
				Target:   &wrapperspb.StringValue{Value: "prometheus"},
				Port:     &wrapperspb.UInt32Value{Value: 8080},
				Path:     &wrapperspb.StringValue{Value: "/metrics"},
				Protocol: &wrapperspb.StringValue{Value: "http"},
			},
		},
	}

	ret, err := convertToClientObject(model.NewClient(client))
	assert.NoError(t, err)

	cop, err := convertToModelClient(ret)
	assert.NoError(t, err)

	cop.Proto().Mtime = nil
	cop.Proto().Ctime = nil

	assert.Equal(t, client, cop.Proto())
}

func createMockClients(total int) []*model.Client {
	ret := make([]*model.Client, 0, total)

	for i := 0; i < 10; i++ {
		ret = append(ret, model.NewClient(&api.Client{
			Host:    &wrapperspb.StringValue{Value: fmt.Sprintf("client-%d", i)},
			Type:    0,
			Version: &wrapperspb.StringValue{Value: fmt.Sprintf("client-%d", i)},
			Location: &api.Location{
				Region: &wrapperspb.StringValue{Value: fmt.Sprintf("client-%d-region", i)},
				Zone:   &wrapperspb.StringValue{Value: fmt.Sprintf("client-%d-zone", i)},
				Campus: &wrapperspb.StringValue{Value: fmt.Sprintf("client-%d-campus", i)},
			},
			Id: &wrapperspb.StringValue{Value: fmt.Sprintf("client-%d", i)},
			Stat: []*api.StatInfo{
				{
					Target:   &wrapperspb.StringValue{Value: "prometheus"},
					Port:     &wrapperspb.UInt32Value{Value: 8080},
					Path:     &wrapperspb.StringValue{Value: "/metrics"},
					Protocol: &wrapperspb.StringValue{Value: "http"},
				},
			},
		}))
	}

	return ret
}

func Test_clientStore_GetMoreClients(t *testing.T) {
	CreateTableDBHandlerAndRun(t, "Test_clientStore_GetMoreClients", func(t *testing.T, handler BoltHandler) {
		cStore := &clientStore{handler: handler}

		total := 10

		mockClients := createMockClients(total)
		err := cStore.BatchAddClients(mockClients)
		assert.NoError(t, err, "batch create clients")

		time.Sleep(time.Second)

		// 首次拉取， mtime 不做处理
		ret, err := cStore.GetMoreClients(time.Now().Add(-1*time.Minute), true)
		assert.NoError(t, err, "get more clients")
		assert.Equal(t, total, len(ret), "get more clients")

		// 非次拉取， mtime 做处理
		ret, err = cStore.GetMoreClients(time.Now().Add(-1*time.Minute), false)
		assert.NoError(t, err, "get more clients")
		assert.Equal(t, total, len(ret), "get more clients")

		// 非次拉取， mtime 做处理
		ret, err = cStore.GetMoreClients(time.Now().Add(10*time.Minute), false)
		assert.NoError(t, err, "get more clients")
		assert.Equal(t, 0, len(ret), "get more clients")
	})
}
