/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http

import (
	"time"

	"github.com/alipay/sofa-mosn/pkg/istio/control"
	"github.com/alipay/sofa-mosn/pkg/istio/utils"
	"github.com/gogo/protobuf/proto"
	"istio.io/api/mixer/v1"
)

type attributesBuilder struct {
	requestContext *control.RequestContext
}

func newAttributesBuilder(requestContext *control.RequestContext) *attributesBuilder {
	return &attributesBuilder{
		requestContext:requestContext,
	}
}

func (b *attributesBuilder) ExtractForwardedAttributes(checkData *CheckData) {
	d, ret := checkData.ExtractIstioAttributes()
	if !ret {
		return
	}
	var attibutes v1.Attributes
	err := proto.Unmarshal([]byte(d), &attibutes)
	if err != nil {
		return
	}

	proto.Merge(&b.requestContext.Attributes, &attibutes)
}

func (b *attributesBuilder) ExtractCheckAttributes(checkData *CheckData) {
	builder := utils.NewAttributesBuilder(&b.requestContext.Attributes)

	srcIp, _, ret := checkData.GetSourceIpPort()
	if ret {
		builder.AddBytes(utils.KOriginIp, []byte(srcIp))
	}

	// TODO: add IsMutualTLS、requested_server_name

	builder.AddTimestamp(utils.KRequestTime, time.Now())

	// TODO: add grpc protocol check
	protocol := "http"
	builder.AddString(utils.KContextProtocol, protocol)

}

func (b *attributesBuilder) ExtractReportAttributes(reportData *ReportData) {
	builder := utils.NewAttributesBuilder(&b.requestContext.Attributes)

	destIp, despPort, err := reportData.GetDestinationIpPort()
	if err == nil {
		if !builder.HasAttribute(utils.KDestinationIp) {
			builder.AddBytes(utils.KDestinationIp, []byte(destIp))
		}
		if !builder.HasAttribute(utils.KDestinationPort) {
			builder.AddInt64(utils.KDestinationIp, int64(despPort))
		}
	}

	headers := reportData.respHeaders
	builder.AddStringMap(utils.KResponseHeaders, headers)

	builder.AddTimestamp(utils.KResponseTime, time.Now())

	reportInfo := reportData.GetReportInfo()

	builder.AddInt64(utils.KRequestBodySize, int64(reportInfo.requestBodySize))
	builder.AddInt64(utils.KResponseBodySize, int64(reportInfo.responseBodySize))
	builder.AddInt64(utils.KRequestTotalSize, int64(reportInfo.requestTotalSize))
	builder.AddInt64(utils.KResponseTotalSize, int64(reportInfo.responseTotalSize))
	builder.AddDuration(utils.KResponseDuration, reportInfo.duration)

	// TODO: add check status code
	builder.AddInt64(utils.KResponseCode, int64(reportInfo.responseCode))

	// TODO: add grpc status report

	// TODO: add response flag
	//builder.AddString()

	// TODO: add rabc info
}

