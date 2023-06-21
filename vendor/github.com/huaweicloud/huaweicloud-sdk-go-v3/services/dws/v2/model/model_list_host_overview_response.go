package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListHostOverviewResponse struct {

	// openApi查询主机概览
	Body           *[]HostOverviewResponse `json:"body,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o ListHostOverviewResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHostOverviewResponse struct{}"
	}

	return strings.Join([]string{"ListHostOverviewResponse", string(data)}, " ")
}
