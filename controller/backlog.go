package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
	"github.com/ffromani/go-todo-app/model"
)

func (ctrl *Controller) BacklogIndex(w http.ResponseWriter, r *http.Request) {
	items, err := ctrl.ld.Filter(func(todo model.Todo) bool {
		return todo.IsOngoing()
	})
	if err != nil {
		sendError(w, http.StatusUnprocessableEntity, err)
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

func (ctrl *Controller) BacklogAssigned(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assignee, ok := vars["assignee"]
	if !ok {
		sendError(w, http.StatusInternalServerError, fmt.Errorf("missing assignee"))
		return
	}
	items, err := ctrl.ld.Filter(func(todo model.Todo) bool {
		return todo.IsOngoing() && todo.Assignee == assignee
	})
	if err != nil {
		sendError(w, http.StatusUnprocessableEntity, err)
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
