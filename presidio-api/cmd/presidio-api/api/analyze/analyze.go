package analyze

import (
	"context"
	"fmt"

	types "github.com/Microsoft/presidio-genproto/golang"
	store "github.com/Microsoft/presidio/presidio-api/cmd/presidio-api/api"
	"github.com/Microsoft/presidio/presidio-api/cmd/presidio-api/api/templates"
)

//Analyze text
func Analyze(ctx context.Context, api *store.API, analyzeAPIRequest *types.AnalyzeApiRequest, project string) ([]*types.AnalyzeResult, error) {

	if analyzeAPIRequest.AnalyzeTemplateId == "" && IsEmptyTemplate(analyzeAPIRequest.AnalyzeTemplate) {
		return nil, fmt.Errorf("Analyze template is missing or empty")
	} else if analyzeAPIRequest.AnalyzeTemplateId != "" && analyzeAPIRequest.AnalyzeTemplateId != nil && !IsEmptyTemplate(analyzeAPIRequest.AnalyzeTemplate) {
		return nil, fmt.Errorf("Both template id and template supplied")			
	} else if analyzeAPIRequest.AnalyzeTemplate == nil {
		analyzeAPIRequest.AnalyzeTemplate = &types.AnalyzeTemplate{}
	}

	err := templates.GetTemplate(api, project, store.Analyze, analyzeAPIRequest.AnalyzeTemplateId, analyzeAPIRequest.AnalyzeTemplate)
	if err != nil {
		return nil, err
	}

	res, err := api.Services.AnalyzeItem(ctx, analyzeAPIRequest.Text, analyzeAPIRequest.AnalyzeTemplate)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("No results")
	}
	return res, err

}

// Breaking down the check for an empty template to accomodate the complex logic of
// checking AllFields and trivial object setup with nil or empty Fields arrays 
// and make it readable.
// Not going deeper than checking for an empty Fields array; a proper check would also
// verify that at least one Fields entry has content.
func IsEmptyTemplate(template *types.AnalyzeTemplate) bool
{
	if template == nil {
		return true
	}
	else if template.AllFields {
		return false
	}
	else {
		return template.Fields == nil || len(template.Fields) == 0
	}
}
