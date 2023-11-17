package base

import (
	"fmt"
	v1 "github.com/xuxant/voyager-operator/api/v1"
	"regexp"
)

func (r *ShipBaseConfigurationReconciler) Validate(ship *v1.Ship) ([]string, error) {
	var message []string
	if msg := r.validateDomainName(); len(msg) > 0 {
		message = append(message, msg...)
	}

	return message, nil
}

func (r *ShipBaseConfigurationReconciler) validateDomainName() []string {
	var message []string

	regex := `^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	if r.Ship.Spec.Domain == "" {
		return message
	}
	match, _ := regexp.MatchString(regex, r.Ship.Spec.Domain)
	if !match {
		message = append(message, fmt.Sprintf("Domain name '%s' is not valid domain name.", r.Ship.Spec.Domain))
	}
	return message
}
