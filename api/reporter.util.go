package api

import "errors"

func (ctx *StateContext) Validate() error {
	if ctx.StateId == "" {
		return errors.New("Context state ID is required.")
	}
	if ctx.StateId == "" {
		return errors.New("State ID is required.")
	}
	return nil
}

func (report *StateReport) Validate() error {
	// Possibly make a null value allowed.
	if report.JsonState == "" && len(report.MsgpackState) == 0 {
		return errors.New("You must include some kind of state in the report.")
	}
	return nil
}

func (req *RecordStateRequest) Validate() error {
	if req.Context == nil {
		return errors.New("Context is required.")
	}
	if req.Report == nil {
		return errors.New("A report is required.")
	}
	if err := req.Context.Validate(); err != nil {
		return err
	}
	if err := req.Report.Validate(); err != nil {
		return err
	}

	return nil
}
