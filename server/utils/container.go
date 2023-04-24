package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func InArr[T primitive.ObjectID | int](arr []T, obj T) bool {
	for _, v := range arr {
		if v == obj {
			return true
		}
	}
	return false
}

func ForEach[T any](arr []T, lambda func(int, T)) {
	for i, v := range arr {
		lambda(i, v)
	}
}
