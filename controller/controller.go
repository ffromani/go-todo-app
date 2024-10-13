package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	apiv1 "github.com/ffromani/go-todo-app/api/v1"
	"github.com/ffromani/go-todo-app/ledger"
	"github.com/ffromani/go-todo-app/middleware"
	"github.com/ffromani/go-todo-app/model"
	"github.com/ffromani/go-todo-app/store"
)

type Controller struct {
	router *mux.Router
	ld     *ledger.Ledger
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func New(ld *ledger.Ledger) http.Handler {
	ctrl := Controller{
		ld:     ld,
		router: mux.NewRouter().StrictSlash(true),
	}
	routes := []Route{
		Route{
			"Index",
			"GET",
			"/",
			ctrl.Index,
		},
		Route{
			"TodoIndex",
			"GET",
			"/todos",
			ctrl.TodoIndex,
		},
		Route{
			"TodoCreate",
			"POST",
			"/todos",
			ctrl.TodoCreate,
		},
		Route{
			"TodoShow",
			"GET",
			"/todos/{todoID}",
			ctrl.TodoShow,
		},
	}

	for _, route := range routes {
		ctrl.router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(middleware.Logger(route.HandlerFunc, route.Name))
	}
	return &ctrl
}

func (ctrl *Controller) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctrl.router.ServeHTTP(w, req)
}

func (ctrl *Controller) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func (ctrl *Controller) TodoIndex(w http.ResponseWriter, r *http.Request) {
	items, err := ctrl.ld.Filter(func(todo model.Todo) bool {
		return true
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		resp := apiv1.Response{
			Status: apiv1.ResponseError,
			Error: &apiv1.Error{
				Code: 422,
				Text: err.Error(),
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
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
		panic(err)
	}
	todo, err := ctrl.ld.Get(store.ID(todoID))
	if err != nil {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		resp := apiv1.Response{
			Status: apiv1.ResponseError,
			Error: &apiv1.Error{
				Code: 404,
				Text: err.Error(),
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	apiTodo := todo.ToAPIv1()
	resp := apiv1.Response{
		Status: apiv1.ResponseSuccess,
		Result: &apiv1.Result{
			Items: []apiv1.Item{
				{
					ID:   apiv1.ID(todoID),
					Todo: &apiTodo,
				},
			},
		},
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos
*/
func (ctrl *Controller) TodoCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	apiTodo, err := apiv1.NewTodoFromJSON(body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		resp := apiv1.Response{
			Status: apiv1.ResponseError,
			Error: &apiv1.Error{
				Code: 422,
				Text: err.Error(),
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	todo := model.NewFromAPITodo(apiTodo)
	log.Printf("API: got object %v", todo)

	todoID, err := ctrl.ld.Set(store.NullID, todo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		resp := apiv1.Response{
			Status: apiv1.ResponseError,
			Error: &apiv1.Error{
				Code: 422,
				Text: err.Error(),
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	resp := apiv1.Response{
		Status: apiv1.ResponseSuccess,
		Result: &apiv1.Result{
			Items: []apiv1.Item{
				{
					ID: apiv1.ID(todoID),
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
