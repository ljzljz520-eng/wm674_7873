package ledger

import (
	"errors"
	"fmt"
)

type Rule struct {
	Name     string
	Minimum  int64
	Maximum  int64
	Currency string
	Enabled  bool
}

func (r Rule) Validate(amount int64, currency string) error {
	if !r.Enabled {
		return nil
	}
	if amount < r.Minimum {
		return fmt.Errorf("amount below %d", r.Minimum)
	}
	if r.Maximum > 0 && amount > r.Maximum {
		return fmt.Errorf("amount above %d", r.Maximum)
	}
	if r.Currency != "" && r.Currency != currency {
		return errors.New("currency mismatch")
	}
	return nil
}

type RuleSet struct{ Rules []Rule }

func DefaultRules() RuleSet {
	return RuleSet{Rules: []Rule{{Name: "minimum", Minimum: 1, Currency: "CNY", Enabled: true}, {Name: "approval-cap", Maximum: 100000000, Currency: "CNY", Enabled: true}}}
}
func (r RuleSet) Check(amount int64, currency string) []error {
	failures := make([]error, 0)
	for _, rule := range r.Rules {
		if err := rule.Validate(amount, currency); err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", rule.Name, err))
		}
	}
	return failures
}
func (r RuleSet) Valid(amount int64, currency string) bool {
	return len(r.Check(amount, currency)) == 0
}
func (r RuleSet) Explain(amount int64, currency string) string {
	failures := r.Check(amount, currency)
	if len(failures) == 0 {
		return "valid"
	}
	message := ""
	for index, failure := range failures {
		if index > 0 {
			message += "; "
		}
		message += failure.Error()
	}
	return message
}
func (r RuleSet) WithCurrency(currency string) RuleSet {
	copyRules := append([]Rule(nil), r.Rules...)
	for index := range copyRules {
		copyRules[index].Currency = currency
	}
	return RuleSet{Rules: copyRules}
}
