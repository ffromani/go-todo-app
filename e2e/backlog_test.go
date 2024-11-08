package e2e_test

// solution:part3

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gotestbootcamp/go-todo-app/model"
)

var _ = Describe("backlog endpoint", Label("backlog"), func() {
	When("todos are added", func() {
		It("should return them", func() {
			todo := model.New("e2e fake todo 1")
			resp, err := httpPost(defaultBaseURL+"/todos", todo.ToAPIv1())
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).To(BeSuccessful())
			Expect(resp).Should(HaveAnyItems())

			objID := resp.Result.Items[0].ID
			resp, err = httpGet(defaultBaseURL + "/backlog")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).To(BeSuccessful())
			Expect(resp).To(HaveItemWithID(objID))
		})
	})
})
