package dto

import "errors"

type DrinkRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Bottle int    `json:"bottle"`
	Cost   int    `json:"cost"`
	Soft   bool   `json:"soft"`
}

func (r *DrinkRequest) Validate() error {
	if r.Bottle < 0 {
		return errors.New("invalid bottle: bottle can't be less than 0")
	}

	if r.Cost < 0 {
		return errors.New("invalid cost: cost can't be less than 0")
	}

	return nil
}
