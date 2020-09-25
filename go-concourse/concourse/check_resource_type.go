package concourse

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/go-concourse/concourse/internal"
)

func (team *team) CheckResourceType(pipelineName string, resourceTypeName string, version atc.Version) (atc.Check, bool, error) {

	params := map[string]string{
		"pipeline_name":      pipelineName,
		"resource_type_name": resourceTypeName,
		"team_name":          team.name,
	}

	var check atc.Check

	jsonBytes, err := json.Marshal(atc.CheckRequestBody{From: version})
	if err != nil {
		return check, false, err
	}

	err = team.connection.Send(internal.Request{
		RequestName: atc.CheckResourceType,
		Params:      params,
		Body:        bytes.NewBuffer(jsonBytes),
		Header:      http.Header{"Content-Type": []string{"application/json"}},
	}, &internal.Response{
		Result: &check,
	})

	switch e := err.(type) {
	case nil:
		return check, true, nil
	case internal.ResourceNotFoundError:
		return check, false, nil
	case internal.UnexpectedResponseError:
		if e.StatusCode == http.StatusInternalServerError {
			return check, false, GenericError{e.Body}
		} else {
			return check, false, err
		}
	default:
		return check, false, err
	}
}
