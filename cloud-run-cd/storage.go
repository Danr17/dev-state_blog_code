package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	errNotFound = errors.New("Item not found, the operation failed")
)

//Memory save tha data in memory
type Memory struct {
	Items []Item
}

func (a *api) listsGoods() ([]Item, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var listItems []Item

	for _, item := range a.db.Items {
		valid := checkValidity(item)
		item.IsValid = valid
		listItems = append(listItems, item)
	}

	return listItems, nil
}

func (a *api) addGood(items ...Item) (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, item := range items {
		for _, i := range a.db.Items {
			if item.ID == i.ID {
				return "item already exists", nil
			}

		}
		addtime := time.Now().Format(layoutRO)
		addtime1, err := time.Parse(layoutRO, addtime)
		if err != nil {
			log.Printf("Can't parse the date, %v", err)
			return "", fmt.Errorf("can't parse the date: %v", err)
		}

		if len(a.db.Items) == 0 {
			item.ID = len(a.db.Items) + 1
		} else {
			lastItem := a.db.Items[len(a.db.Items)-1]
			item.ID = lastItem.ID + 1
		}

		item.Created = timestamp{addtime1}

		a.db.Items = append(a.db.Items, item)
	}

	return "Items has been added to database", nil
}

func (a *api) openState(id int, status bool) (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var foundIndex int
	var found bool
	for i, item := range a.db.Items {
		if id == item.ID {
			found = true
			foundIndex = i
			break
		}
	}
	if !found {
		return "", errNotFound
	}

	opentimeS := time.Now().Format(layoutRO)
	opentimeT, err := time.Parse(layoutRO, opentimeS)
	if err != nil {
		log.Printf("Can't parse the date, %v", err)
		return "", fmt.Errorf("can't parse the date: %v", err)
	}
	a.db.Items[foundIndex].IsOpen = status
	a.db.Items[foundIndex].Opened = timestamp{opentimeT}

	return fmt.Sprintf("Item with id %d has been opened", a.db.Items[foundIndex].ID), nil
}

func (a *api) delGood(id int) (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var foundIndex int
	var found bool
	for i, item := range a.db.Items {
		if id == item.ID {
			found = true
			foundIndex = i
			break
		}
	}

	if !found {
		return "", errNotFound
	}

	a.db.Items = removeIndex(a.db.Items, foundIndex)
	return fmt.Sprintf("Item id %d has been deleted", id), nil

}

func removeIndex(s []Item, index int) []Item {
	return append(s[:index], s[index+1:]...)
}

func checkValidity(i Item) bool {
	t := time.Now()
	i.IsValid = true
	if t.Sub(i.ExpDate.Time) > 0 {
		i.IsValid = false
	}

	if i.IsOpen {
		//if i.Opened.Time.Add(time.Duration(int64(time.Duration(i.ExpOpen * 24).Hours()))).Before(t) {
		if t.Sub(i.Opened.Time.AddDate(0, 0, i.ExpOpen)) > 0 {
			i.IsValid = false
		}
	}

	return i.IsValid
}
