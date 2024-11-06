package e2e_test

// solution:part3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bsm/gomega/gcustom"
	"github.com/bsm/gomega/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
	"github.com/gotestbootcamp/go-todo-app/model"
)

const (
	defaultBaseURL = "http://localhost:8181"
)

var _ = Describe("backlog endpoint", func() {
	When("todos are added", func() {
		It("should return them", func() {
			todo := model.Todo{
				Title:       "e2e fake todo 1",
				Description: "bogus todo for e2e testing",
			}
			data, err := todo.ToAPIv1().ToJSON()
			Expect(err).ToNot(HaveOccurred())
			Expect(data).ToNot(BeEmpty())

			res, err := http.Post(defaultBaseURL+"/todos", "application/json", bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())

			var apiCreateResp apiv1.Response
			Expect(json.NewDecoder(res.Body).Decode(&apiCreateResp)).To(Succeed())

			Expect(apiCreateResp).To(BeSuccessful())
			Expect(len(apiCreateResp.Result.Items)).To(Equal(1))

			objID := apiCreateResp.Result.Items[0].ID

			res, err = http.Get(defaultBaseURL + "/backlog")
			Expect(err).ToNot(HaveOccurred())

			var apiBacklogResp apiv1.Response
			Expect(json.NewDecoder(res.Body).Decode(&apiBacklogResp)).To(Succeed())
			Expect(apiBacklogResp).To(BeSuccessful())
			Expect(apiBacklogResp).To(HaveItemWithID(objID))
		})
	})
})

func BeSuccessful() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(actual apiv1.Response) (bool, error) {
			return actual.Status == apiv1.ResponseSuccess && actual.Error == nil, nil
		},
	).WithTemplate("API v1 response should be succesfull and not have error")
}

func HaveItemWithID(objID apiv1.ID) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(actual apiv1.Response) (bool, error) {
			for _, item := range actual.Result.Items {
				if item.ID == objID {
					return true, nil
				}
			}
			return false, nil
		},
	).WithTemplate("API v1 response should include {{.Data}}").WithTemplateData(objID)
}

func toJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("<JSON marshal err=%v>", err)
	}
	return string(data)
}
