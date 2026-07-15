package urn

import (
	"fmt"
	"strings"
)

const Prefix = "urn:cloud-terminal"

func New(tenant, resourceType, resourceID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", Prefix, tenant, resourceType, resourceID)
}

type Info struct {
	Tenant       string
	ResourceType string
	ResourceID   string
}

func Parse(urn string) (*Info, error) {
	parts := strings.SplitN(urn, ":", 5)
	if len(parts) != 5 || parts[0] != "urn" || parts[1] != "cloud-terminal" {
		return nil, fmt.Errorf("invalid URN: %s", urn)
	}
	return &Info{
		Tenant:       parts[2],
		ResourceType: parts[3],
		ResourceID:   parts[4],
	}, nil
}

func Match(pattern, urn string) bool {
	pParts := strings.SplitN(pattern, ":", 5)
	uParts := strings.SplitN(urn, ":", 5)
	if len(pParts) != 5 || len(uParts) != 5 {
		return false
	}
	for i := 0; i < 5; i++ {
		if pParts[i] == "*" {
			continue
		}
		if pParts[i] != uParts[i] {
			return false
		}
	}
	return true
}
