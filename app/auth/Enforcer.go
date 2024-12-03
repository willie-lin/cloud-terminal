package auth

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

type Enforcer struct {
	enforcer *casbin.Enforcer
}

func NewEnforcer(modelPath, policyPath string) (*Enforcer, error) {
	e, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}
	return &Enforcer{enforcer: e}, nil
}

func (e *Enforcer) AddPolicy(sub, dom, obj, act string) error {
	_, err := e.enforcer.AddPolicy(sub, dom, obj, act)
	return err
}

func (e *Enforcer) AddRoleForUser(user, role, domain string) error {
	_, err := e.enforcer.AddGroupingPolicy(user, role, domain)
	return err
}

func (e *Enforcer) Enforce(sub, dom, obj, act string) (bool, error) {
	return e.enforcer.Enforce(sub, dom, obj, act)
}
