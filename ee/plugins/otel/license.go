// Copyright 2025 samber.
//
// Licensed as an Enterprise License (the "License"); you may not use
// this file except in compliance with the License. You may obtain
// a copy of the License at:
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.ee.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package rootel

import (
	"context"

	"github.com/samber/ro"
	rolicense "github.com/samber/ro/ee/pkg/license"
)

var bypassLicenseCheck = false

func isOtelEnabled() bool {
	return bypassLicenseCheck || rolicense.IsEnterpriseEnabled()
}

func checkLicenseAndPipe[First any, Last any](
	collector *otelCollector,
	source ro.Observable[First],
	stdPipe func(ro.Observable[First]) ro.Observable[Last],
	instrumentedPipe func(ro.Observable[First]) ro.Observable[Last],
) ro.Observable[Last] {
	return ro.NewUnsafeObservableWithContext(func(subscriberCtx context.Context, destination ro.Observer[Last]) ro.Teardown {
		var p func(ro.Observable[First]) ro.Observable[Last]

		if isOtelEnabled() {
			p = wrapPipeWithObservability(collector, instrumentedPipe)
		} else {
			p = stdPipe
		}

		sub := p(source).SubscribeWithContext(subscriberCtx, destination)
		return sub.Unsubscribe
	})
}
