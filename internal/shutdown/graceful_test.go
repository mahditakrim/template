package shutdown

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestGraceful(t *testing.T) {

	type args struct {
		f       func() error
		sigChan chan os.Signal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "proper",
			args: args{
				f:       func() error { return nil },
				sigChan: make(chan os.Signal, 1),
			},
			wantErr: false,
		},
		{
			name: "nil func",
			args: args{
				f:       nil,
				sigChan: make(chan os.Signal, 1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				err := Graceful(tt.args.f, tt.args.sigChan)
				if (err != nil) != tt.wantErr {
					t.Errorf("Graceful() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			time.Sleep(time.Second)
			tt.args.sigChan <- syscall.SIGINT
		})
	}
}
