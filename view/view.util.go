package view

import (
	"errors"
)

func (report *StateReport) Validate() error {
	// Possibly make a null value allowed.
	if report.JsonState == "" {
		return errors.New("You must include some kind of state in the report.")
	}
	return nil
}
func (req *StateQuery) Validate() error {
	return nil
}

func (req *GetStateRequest) Validate(requireHost bool) error {
	if err := req.Context.Validate(requireHost); err != nil {
		return err
	}
	if req.Query == nil {
		req.Query = &StateQuery{}
	}
	if err := req.Query.Validate(); err != nil {
		return err
	}
	return nil
}

func (ctx *StateContext) Validate(requireHost bool) error {
	if requireHost && ctx.HostIdentifier == "" {
		return errors.New("Host identifier must be specified.")
	}
	if ctx.StateId == "" {
		return errors.New("Context state ID is required.")
	}
	if ctx.StateId == "" {
		return errors.New("State ID is required.")
	}
	return nil
}

func (r *StateHistoryQuery) Validate() error {
	if r == nil {
		return errors.New("Query is required.")
	}
	if r.EndTime != 0 && r.Tail {
		return errors.New("Cannot tail unless EndTime = 0 (indicating latest).")
	}
	if r.BeginTime > r.EndTime {
		return errors.New("BeginTime cannot be greater than EndTime.")
	}
	return nil
}
