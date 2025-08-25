package domain

import (
	"fmt"
	"math/rand"
	"strings"
)

type ResourceType string

const (
	ResourceTypeQuestion   ResourceType = "qst"
	ResourceTypeChest      ResourceType = "cst"
	ResourceTypeComponent  ResourceType = "cmp"
	ResourceTypeAnswer     ResourceType = "ans"
	ResourceTypeCorrection ResourceType = "crt"
)

func NewID(resourceType ResourceType) string {
	id := rand.Int63()
	return fmt.Sprintf("%s_%X", resourceType, id)
}

func IdHasType(id string, resourceType ResourceType) bool {
	return strings.HasPrefix(id, string(resourceType)+"_")
}
