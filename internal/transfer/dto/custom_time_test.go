package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnsureTruncatedCorrectly(t *testing.T) {
	t.Run("time is truncated to 6 digits", func(t *testing.T) {
		timeNow := time.Now()

		// yes this is a reimplementation, but this must be precise to no more than 6 digits
		// as postgres cannot handle time more sensitive than that
		// is this a leakage of repository limitation to app tier?
		// Yes
		// Can this be fixed?
		// Also yes, add additional column to postgres that contains the remaining digits
		// But I thought of that too late and now here we are
		// This is how tech debt is built :)
		assert.Equal(t, timeNow.UTC().Truncate(time.Microsecond), NewCustomTime(timeNow).Time)
	})
}
