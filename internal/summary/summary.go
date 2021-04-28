package summary

import (
	audit "github.com/citihub/probr-sdk/audit"
)

// State should be set in the pack's runtime via audit.NewSummaryState
var State audit.SummaryState
