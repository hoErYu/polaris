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

package v2

import (
	"github.com/polarismesh/polaris-server/namespace"
	"github.com/polarismesh/polaris-server/service"
	"github.com/polarismesh/polaris-server/service/healthcheck"
)

// HTTPServerV2
type HTTPServerV2 struct {
	namespaceServer   namespace.NamespaceOperateServer
	namingServer      service.DiscoverServer
	healthCheckServer *healthcheck.Server
}

// NewV2Server
func NewV2Server(
	namespaceServer namespace.NamespaceOperateServer,
	namingServer service.DiscoverServer,
	healthCheckServer *healthcheck.Server) *HTTPServerV2 {

	return &HTTPServerV2{
		namespaceServer:   namespaceServer,
		namingServer:      namingServer,
		healthCheckServer: healthCheckServer,
	}
}
