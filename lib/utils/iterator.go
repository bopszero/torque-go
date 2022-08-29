package utils

import (
	"reflect"
)

type (
	IteratorLoader       func(lastItem interface{}) (items interface{}, err error)
	IteratorGetNext      func(value interface{}) (IteratorGetNext, IteratorState)
	IteratorGetNextChunk func(value interface{}) (IteratorGetNextChunk, IteratorState)
)

type IteratorState struct {
	err       error
	exhausted bool
}

func (this *IteratorState) OK() bool {
	return !this.exhausted && this.err == nil
}

func (this *IteratorState) Exhausted() bool {
	return this.exhausted
}

func (this *IteratorState) Error() error {
	return this.err
}

type ObjectIterator struct {
	loader     IteratorLoader
	latestItem interface{}
}

func NewObjectIterator(loader IteratorLoader) ObjectIterator {
	return ObjectIterator{
		loader: loader,
	}
}

func (this *ObjectIterator) GetNext(value interface{}) (IteratorGetNext, IteratorState) {
	return this.pick(value, reflect.ValueOf([]struct{}(nil)))
}

func (this *ObjectIterator) pick(value interface{}, items reflect.Value) (IteratorGetNext, IteratorState) {
	if items.Len() == 0 {
		loadedItems, err := this.load()
		if err != nil {
			return nil, IteratorState{err: err}
		}
		if loadedItems == nil {
			return nil, IteratorState{exhausted: true}
		}
		items = reflect.ValueOf(loadedItems)
		if items.Len() == 0 {
			return nil, IteratorState{exhausted: true}
		}
	}
	var (
		pickItem    = items.Index(0).Interface()
		pickValue   = reflect.Indirect(reflect.ValueOf(pickItem))
		targetValue = reflect.Indirect(reflect.ValueOf(value))
	)
	if pickValue.Type() != targetValue.Type() {
		err := IssueErrorf(
			"iterator cannot set type `%v` to target type `%v`",
			pickValue.Type(), targetValue.Type())
		return nil, IteratorState{err: err}
	}
	reflect.Indirect(reflect.ValueOf(value)).Set(pickValue)

	this.latestItem = pickItem
	items = items.Slice(1, items.Len())
	getNext := func(value interface{}) (IteratorGetNext, IteratorState) {
		return this.pick(value, items)
	}
	return getNext, IteratorState{}
}

func (this *ObjectIterator) load() (interface{}, error) {
	items, err := this.loader(this.latestItem)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return nil, nil
	}
	itemsType := reflect.TypeOf(items)
	switch itemsType.Kind() {
	case reflect.Slice, reflect.Array:
		break
	default:
		err := IssueErrorf(
			"iterator expects loader items should be Slice/Array but receive `%v`",
			itemsType)
		return nil, err
	}
	return items, nil
}

type ChunkIterator struct {
	loader     IteratorLoader
	size       int
	latestItem interface{}
}

func NewChunkIterator(size int, loader IteratorLoader) ChunkIterator {
	return ChunkIterator{
		loader: loader,
		size:   size,
	}
}

func (this *ChunkIterator) GetNext(values interface{}) (IteratorGetNextChunk, IteratorState) {
	return this.pick(values, reflect.ValueOf([]struct{}(nil)))
}

func (this *ChunkIterator) pick(values interface{}, items reflect.Value) (IteratorGetNextChunk, IteratorState) {
	if items.Len() < this.size {
		loadedItems, err := this.load()
		if err != nil {
			return nil, IteratorState{err: err}
		}
		if loadedItems != nil {
			loadedItemsValue := reflect.ValueOf(loadedItems)
			if items.Len() == 0 {
				items = loadedItemsValue
			} else {
				items = reflect.AppendSlice(items, loadedItemsValue)
			}
		}
		if items.Len() == 0 {
			return nil, IteratorState{exhausted: true}
		}
	}
	pickSize := this.size
	if items.Len() < this.size {
		pickSize = items.Len()
	}
	var (
		pickItems      = items.Slice(0, pickSize)
		targetItems    = reflect.Indirect(reflect.ValueOf(values))
		newTargetItems = reflect.MakeSlice(targetItems.Type(), 0, pickItems.Len())
	)
	targetItems.Set(reflect.AppendSlice(newTargetItems, pickItems))

	this.latestItem = pickItems.Index(pickItems.Len() - 1).Interface()
	items = items.Slice(pickSize, items.Len())
	getNext := func(values interface{}) (IteratorGetNextChunk, IteratorState) {
		return this.pick(values, items)
	}
	return getNext, IteratorState{}
}

func (this *ChunkIterator) load() (interface{}, error) {
	items, err := this.loader(this.latestItem)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return nil, nil
	}
	itemsType := reflect.TypeOf(items)
	switch itemsType.Kind() {
	case reflect.Slice, reflect.Array:
		break
	default:
		err := IssueErrorf(
			"iterator expects loader items should be Slice/Array but receive `%v`",
			itemsType)
		return nil, err
	}
	return items, nil
}
