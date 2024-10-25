package fake

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ffromani/go-todo-app/store"
	"github.com/stretchr/testify/assert"
)

func TestNewEmpty(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)
	items, err := st.LoadAll()
	assert.NoError(t, err)
	assert.Empty(t, items)
}

func TestLoadAllLoad(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)
	num := 0
	count := 5 // random number > 1
	st.Generate = func() (store.Item, bool, error) {
		done := (num == count)
		item := store.Item{
			ID:   store.ID(fmt.Sprintf("%d", num)),
			Blob: store.Blob(fmt.Sprintf("data#%d", num)),
		}
		num += 1
		return item, done, nil
	}
	items, err := st.LoadAll()
	assert.Equal(t, len(items), count)
}

func TestClose(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)
	err = st.Close()
	assert.NoError(t, err)
}

func TestCloseWithError(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)
	expErr := errors.New("injected error")
	st.Error = expErr
	err = st.Close()
	assert.ErrorIs(t, err, expErr)
}

func TestCreateError(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	expErr := errors.New("injected create error")
	st.Error = expErr

	val := "foobar"
	err = st.Create("1", store.Blob(val))
	assert.ErrorIs(t, err, expErr)
}

func TestCreateLoad(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	val := "foobar"
	err = st.Create("1", store.Blob(val))
	assert.NoError(t, err)

	blob, err := st.Load("1")
	assert.Equal(t, string(blob), val, "retrieved object different from inserted")
}

func TestLoadWithError(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	val := "foobar"
	err = st.Create("1", store.Blob(val))
	assert.NoError(t, err)

	// without this error injection, Load() will succeed
	expErr := errors.New("injected load error")
	st.Error = expErr

	_, err = st.Load("1")
	assert.ErrorIs(t, err, expErr)
}

func TestLoadFromEmpty(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	_, err = st.Load("999")
	assert.ErrorIs(t, err, store.ErrNotFound{ID: "999"})
}

func TestSaveFromEmpty(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	val := "foobar"
	err = st.Save("999", store.Blob(val))
	assert.ErrorIs(t, err, store.ErrNotFound{ID: "999"})
}

func TestDeleteWithError(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	expErr := errors.New("injected delete error")
	st.Error = expErr

	err = st.Delete(store.ID("999"))
	assert.ErrorIs(t, err, expErr)
}

func TestDeleteFromEmpty(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	err = st.Delete("999")
	assert.ErrorIs(t, err, store.ErrNotFound{ID: "999"})
}

func TestCreateSaveLoad(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	val := "foobar"
	err = st.Create("123", store.Blob(val))
	assert.NoError(t, err)

	val2 := "fizzbuzz"
	err = st.Save("123", store.Blob(val2))
	assert.NoError(t, err)

	blob, err := st.Load("123")
	assert.Equal(t, string(blob), val2, "retrieved object different from insterted")
}

func TestCreateDeleteLoad(t *testing.T) {
	st, err := NewMem()
	assert.NoError(t, err)

	id := store.ID("543")
	val := "foobar"
	err = st.Create(id, store.Blob(val))
	assert.NoError(t, err)

	err = st.Delete(id)
	assert.NoError(t, err)

	_, err = st.Load(id)
	assert.ErrorIs(t, err, store.ErrNotFound{ID: id})
}
