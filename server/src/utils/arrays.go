package utils

import "slices"

func MapWithoutDuplicate[T any, U comparable](array []T, mappingFunc func(param T) U) []U {
	var mappedSet []U
	for i := range array {
		result := mappingFunc(array[i])
		if !slices.Contains(mappedSet, result) {
			mappedSet = append(mappedSet, result)
		}
	}
	return mappedSet
}

func Map[T any, U any](array []T, mappingFunc func(param T) U) []U {
	var mappedArray []U
	for i := range array {
		mappedArray = append(mappedArray, mappingFunc(array[i]))
	}
	return mappedArray
}

func Filter[T any](array []T, predicate func(param T) bool) []T {
	var filteredArray []T
	for i := range array {
		if predicate(array[i]) {
			filteredArray = append(filteredArray, array[i])
		}
	}
	return filteredArray
}
