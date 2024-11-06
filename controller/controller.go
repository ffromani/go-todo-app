package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	apiv1 "github.com/gotestbootcamp/go-todo-app/api/v1"
	"github.com/gotestbootcamp/go-todo-app/ledger"
	"github.com/gotestbootcamp/go-todo-app/middleware"
	"github.com/gotestbootcamp/go-todo-app/uuid"
)

type Controller struct {
	router  *mux.Router
	ld      *ledger.Ledger
	uuidGen uuid.UUIDGenerator
}

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

func New(ld *ledger.Ledger) http.Handler {
	ctrl := Controller{
		ld:      ld,
		uuidGen: uuid.New(),
		router:  mux.NewRouter().StrictSlash(true),
	}
	routes := []Route{
		Route{
			Name:    "backlog.index",
			Method:  "GET",
			Pattern: "/backlog",
			Handler: ctrl.BacklogIndex,
		},
		Route{
			Name:    "backlog.assigned",
			Method:  "GET",
			Pattern: "/backlog/{assignee}",
			Handler: ctrl.BacklogAssigned,
		},
		Route{
			Name:    "completed.index",
			Method:  "GET",
			Pattern: "/completed",
			Handler: ctrl.CompletedIndex,
		},
		Route{
			Name:    "completed.byassignee",
			Method:  "GET",
			Pattern: "/completed/{assignee}",
			Handler: ctrl.CompletedAssigned,
		},
		Route{
			Name:    "todo.index",
			Method:  "GET",
			Pattern: "/todos",
			Handler: ctrl.TodoIndex,
		},
		Route{
			Name:    "todo.create",
			Method:  "POST",
			Pattern: "/todos",
			Handler: ctrl.TodoCreate,
		},
		Route{
			Name:    "todo.show",
			Method:  "GET",
			Pattern: "/todos/{todoID}",
			Handler: ctrl.TodoShow,
		},
		// PUT is defined to assume idempotency, so if you PUT an object twice, it should have no additional effect.
		Route{
			Name:    "todo.update",
			Method:  "PUT",
			Pattern: "/todos/{todoID}",
			Handler: ctrl.TodoUpdate,
		},
		// you can complete a TODO just once
		Route{
			Name:    "todo.complete",
			Method:  "POST",
			Pattern: "/todos/{todoID}/complete",
			Handler: ctrl.TodoComplete,
		},
		// you can delete a TODO just once
		Route{
			Name:    "todo.delete",
			Method:  "POST",
			Pattern: "/todos/{todoID}/delete",
			Handler: ctrl.TodoDelete,
		},
		Route{
			Name:    "todo.merge",
			Method:  "POST",
			Pattern: "/todomerge/{todoID1}/{todoID2}",
			Handler: ctrl.TodoMerge,
		},
	}

	for _, route := range routes {
		ctrl.router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(middleware.Logger(route.Handler, route.Name))
		log.Printf("API: method: %-8s route: %s", route.Method, route.Pattern)
	}
	return &ctrl
}

func (ctrl *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctrl.router.ServeHTTP(w, req)
}

func sendError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	resp := apiv1.Response{
		Status: apiv1.ResponseError,
		Error: &apiv1.Error{
			Code: code,
			Text: err.Error(),
		},
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func sendItem(w http.ResponseWriter, id apiv1.ID, todo *apiv1.Todo) {
	resp := apiv1.Response{
		Status: apiv1.ResponseSuccess,
		Result: &apiv1.Result{
			Items: []apiv1.Item{
				{
					ID:   id,
					Todo: todo,
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
