package date

import (
	"testing"
	"time"
)

func TestWithinDuration(t *testing.T) {
	type args struct {
		expected time.Time
		actual   time.Time
		delta    time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is within duration",
			args: args{
				expected: time.Date(2021, 12, 21, 0, 0, 0, 0, time.UTC),
				actual:   time.Date(2021, 12, 21, 0, 0, 1, 0, time.UTC),
				delta:    time.Millisecond * 1500,
			},
			want: true,
		},
		{
			name: "not within duration",
			args: args{
				expected: time.Date(2021, 12, 21, 0, 0, 0, 0, time.UTC),
				actual:   time.Date(2021, 12, 21, 0, 0, 1, 0, time.UTC),
				delta:    time.Millisecond * 900,
			},
			want: false,
		},
		{
			name: "not inclusive start",
			args: args{
				expected: time.Date(2021, 12, 21, 0, 0, 0, 0, time.UTC),
				actual:   time.Date(2021, 12, 21, 0, 0, 1, 0, time.UTC),
				delta:    time.Second,
			},
			want: false,
		},
		{
			name: "not inclusive end",
			args: args{
				expected: time.Date(2021, 12, 21, 0, 0, 1, 0, time.UTC),
				actual:   time.Date(2021, 12, 21, 0, 0, 0, 0, time.UTC),
				delta:    time.Second,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithinDuration(tt.args.expected, tt.args.actual, tt.args.delta); got != tt.want {
				t.Errorf("WithinDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
