package e2e_test

// solution:part3

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gotestbootcamp/go-todo-app/model"
)

var _ = Describe("todo lifecycle", Label("lifecycle"), func() {
	When("todos are added, assigned and completed", func() {
		It("should return them after each step", func() {
			todo := model.New("e2e fake todo 2")
			resp, err := httpPost(defaultBaseURL+"/todos", todo.ToAPIv1())
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).ToNot(BeNil())
			Expect(resp).To(BeSuccessful())
			Expect(resp).Should(HaveAnyItems())

			objID := resp.Result.Items[0].ID

			resp, err = httpGet(defaultBaseURL + "/completed")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).To(BeSuccessful())
			Expect(resp).ShouldNot(HaveAnyItems())

			userName := "JacquelineDoe"
			err = todo.Assign(userName)
			Expect(err).ToNot(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "updated todo: %s\n", toJSON(todo))

			resp, err = httpPut(defaultBaseURL+"/todos/"+string(objID), todo.ToAPIv1())
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).ToNot(BeNil())
			Expect(resp).To(BeSuccessful())

			resp, err = httpGet(defaultBaseURL + "/backlog/" + userName)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).To(BeSuccessful())
			Expect(resp).To(HaveItemWithID(objID))

			err = todo.Complete()
			Expect(err).ToNot(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "updated todo: %s\n", toJSON(todo))
			resp, err = httpPost(defaultBaseURL+"/todos/"+string(objID)+"/complete", todo.ToAPIv1())
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).ToNot(BeNil())
			Expect(resp).To(BeSuccessful())

			resp, err = httpGet(defaultBaseURL + "/completed/" + userName)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).To(BeSuccessful())
			Expect(resp).To(HaveItemWithID(objID))
		})
	})
})
