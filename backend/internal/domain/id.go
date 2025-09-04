package domain

import (
	"fmt"
	"math/rand"
	"strings"
)

type ResourceType string

const (
	ResourceTypeBook       ResourceType = "bok"
	ResourceTypeQuestion   ResourceType = "qst"
	ResourceTypeTreasure   ResourceType = "trs"
	ResourceTypeCorrection ResourceType = "crt"
)

func NewID(resourceType ResourceType) string {
	id := rand.Int63()
	return fmt.Sprintf("%s_%X", resourceType, id)
}

func IdHasType(id string, resourceType ResourceType) bool {
	return strings.HasPrefix(id, string(resourceType)+"_")
}
