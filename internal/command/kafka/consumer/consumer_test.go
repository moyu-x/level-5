package consumer_test

import (
	"testing"

	"github.com/expr-lang/expr"
)

func TestExpr(t *testing.T) {
	_, err := expr.Compile("")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
