package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
	"github.com/ffromani/go-todo-app/model"
	"github.com/ffromani/go-todo-app/store"
)

func (ctrl *Controller) TodoIndex(w http.ResponseWriter, r *http.Request) {
	items, err := ctrl.ld.Filter(func(todo model.Todo) bool {
		return true
	})
	if err != nil {
		sendError(w, 422, err)
		return
	}

	resp := apiv1.Response{
		Status: apiv1.ResponseSuccess,
		Result: &apiv1.Result{
			Items: items.ToAPIv1(),
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func (ctrl *Controller) TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todoID int
	var err error
	if todoID, err = strconv.Atoi(vars["todoID"]); err != nil {
		sendError(w, 422, err)
		return
	}
	todo, err := ctrl.ld.Get(store.ID(todoID))
	if err != nil {
		sendError(w, 404, err)
		return
	}

	apiTodo := todo.ToAPIv1()
	sendItem(w, apiv1.ID(todoID), &apiTodo)
}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos
*/
func (ctrl *Controller) TodoCreate(w http.ResponseWriter, r *http.Request) {
	apiTodo, code, err := todoFromRequest(r)
	if err != nil {
		sendError(w, code, err)
		return
	}

	todo := model.NewFromAPIv1(apiTodo)
	log.Printf("API: got object %v", todo)

	todoID, err := ctrl.ld.Set(store.NullID, todo)
	if err != nil {
		sendError(w, 422, err)
		return
	}

	sendItem(w, apiv1.ID(todoID), nil)
}

func (ctrl *Controller) TodoUpdate(w http.ResponseWriter, r *http.Request) {
	apiTodo, code, err := todoFromRequest(r)
	if err != nil {
		sendError(w, code, err)
		return
	}

	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["todoID"])
	if err != nil {
		sendError(w, 500, err)
		return
	}
	todo, err := ctrl.ld.Get(store.ID(todoID))
	if err != nil {
		sendError(w, 404, err)
		return
	}
	log.Printf("API: got object %v", todoID)

	if err := todo.Describe(apiTodo.Description); err != nil {
		sendError(w, 422, err)
		return
	}
	if err := todo.Assign(apiTodo.Assignee); err != nil {
		sendError(w, 422, err)
		return
	}

	log.Printf("API: updated object %v as: %q", todoID, todo)

	_, err = ctrl.ld.Set(store.ID(todoID), todo)
	if err != nil {
		sendError(w, 422, err)
		return
	}

	resTodo := todo.ToAPIv1()
	sendItem(w, apiv1.ID(todoID), &resTodo)
}

func (ctrl *Controller) TodoComplete(w http.ResponseWriter, r *http.Request) {
	_, code, err := todoFromRequest(r)
	if err != nil {
		sendError(w, code, err)
		return
	}

	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["todoID"])
	if err != nil {
		sendError(w, 500, err)
		return
	}
	todo, err := ctrl.ld.Get(store.ID(todoID))
	if err != nil {
		sendError(w, 404, err)
		return
	}
	log.Printf("API: got object %v", todoID)

	if err := todo.Complete(); err != nil {
		sendError(w, 422, err)
		return
	}

	log.Printf("API: completed object %v as: %q", todoID, todo)

	_, err = ctrl.ld.Set(store.ID(todoID), todo)
	if err != nil {
		sendError(w, 422, err)
		return
	}

	resTodo := todo.ToAPIv1()
	sendItem(w, apiv1.ID(todoID), &resTodo)
}

func (ctrl *Controller) TodoDelete(w http.ResponseWriter, r *http.Request) {
	_, code, err := todoFromRequest(r)
	if err != nil {
		sendError(w, code, err)
		return
	}

	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["todoID"])
	if err != nil {
		sendError(w, 500, err)
		return
	}
	todo, err := ctrl.ld.Get(store.ID(todoID))
	if err != nil {
		sendError(w, 404, err)
		return
	}
	log.Printf("API: got object %v", todoID)

	if err := todo.Delete(); err != nil {
		sendError(w, 422, err)
		return
	}

	log.Printf("API: deleted object %v as: %q", todoID, todo)

	_, err = ctrl.ld.Set(store.ID(todoID), todo)
	if err != nil {
		sendError(w, 422, err)
		return
	}

	resTodo := todo.ToAPIv1()
	sendItem(w, apiv1.ID(todoID), &resTodo)
}

func (ctrl *Controller) TodoMerge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id1, err := strconv.Atoi(vars["todoID1"])
	if err != nil {
		sendError(w, 500, err)
		return
	}

	id2, err := strconv.Atoi(vars["todoID2"])
	if err != nil {
		sendError(w, 500, err)
		return
	}
	todo1, err := ctrl.ld.Get(store.ID(id1))
	if err != nil {
		sendError(w, 404, err)
		return
	}
	todo2, err := ctrl.ld.Get(store.ID(id2))
	if err != nil {
		sendError(w, 404, err)
		return
	}
	log.Printf("API: got objects %v - %v", todo1, todo2)

	merged, err := model.Merge(todo1, todo2)
	if err != nil {
		sendError(w, 422, err)
		return
	}

	err = ctrl.ld.Delete(store.ID(id1))
	if err != nil {
		sendError(w, 422, err)
		return
	}
	err = ctrl.ld.Delete(store.ID(id2))
	if err != nil {
		sendError(w, 422, err)
		return
	}

	mergedID, err := ctrl.ld.Set(store.NullID, merged)
	if err != nil {
		sendError(w, 422, err)
		return
	}

	resTodo := merged.ToAPIv1()
	sendItem(w, apiv1.ID(mergedID), &resTodo)
}

func todoFromRequest(r *http.Request) (apiv1.Todo, int, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return apiv1.Todo{}, 500, err
	}
	if err := r.Body.Close(); err != nil {
		return apiv1.Todo{}, 500, err
	}
	apiTodo, err := apiv1.NewTodoFromJSON(body)
	if err != nil {
		return apiv1.Todo{}, 500, err
	}
	return apiTodo, 0, nil
}
