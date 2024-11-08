package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bsm/gomega/gcustom"
	"github.com/bsm/gomega/types"
	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	defaultBaseURL = "http://localhost:8181"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

func BeSuccessful() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(actual *apiv1.Response) (bool, error) {
			return actual.Status == apiv1.ResponseSuccess && actual.Error == nil, nil
		},
	).WithTemplate("API v1 response should be succesfull and not have error")
}

func HaveItemWithID(objID apiv1.ID) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(actual *apiv1.Response) (bool, error) {
			for _, item := range actual.Result.Items {
				if item.ID == objID {
					return true, nil
				}
			}
			return false, nil
		},
	).WithTemplate("API v1 response should include {{.Data}}").WithTemplateData(objID)
}

func HaveAnyItems() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(actual *apiv1.Response) (bool, error) {
			return len(actual.Result.Items) > 0, nil
		},
	).WithTemplate("API v1 response should include any items")
}

func toJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("<JSON marshal err=%v>", err)
	}
	return string(data)
}

func httpPost(url string, obj apiv1.Todo) (*apiv1.Response, error) {
	data, err := obj.ToJSON()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var apiResp apiv1.Response
	err = json.NewDecoder(res.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	return &apiResp, nil
}

func httpPut(url string, obj apiv1.Todo) (*apiv1.Response, error) {
	data, err := obj.ToJSON()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var apiResp apiv1.Response
	err = json.NewDecoder(res.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	return &apiResp, nil
}

func httpGet(url string) (*apiv1.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var apiResp apiv1.Response
	err = json.NewDecoder(res.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(GinkgoWriter, "apiResponse: %v\n", toJSON(apiResp))
	return &apiResp, nil
}
