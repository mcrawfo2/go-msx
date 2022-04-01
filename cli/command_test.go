// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cli

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestAddCommand(t *testing.T) {
	cmdFunc := func(args []string) error {
		return nil
	}

	type args struct {
		path    string
		brief   string
		cmdFunc CommandFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Primary",
			args: args{
				path:    "primary",
				brief:   "primary command",
				cmdFunc: cmdFunc,
			},
			wantErr: false,
		},
		{
			name: "Secondary",
			args: args{
				path:    "primary secondary",
				brief:   "secondary command",
				cmdFunc: cmdFunc,
			},
			wantErr: false,
		},
		{
			name: "MissingSecondary",
			args: args{
				path:    "missing secondary",
				brief:   "missing secondary command",
				cmdFunc: cmdFunc,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddCommand(tt.args.path, tt.args.brief, tt.args.cmdFunc)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AddCommand() error = %v, wantErr %v", err, tt.wantErr)

				}
				assert.Nil(t, got)
				return
			} else if got == nil {
				assert.NotNil(t, got)
				return
			}

			if strings.TrimSpace(got.CommandPath()) != tt.args.path {
				t.Errorf("Path got = %q, want %q", strings.TrimSpace(got.CommandPath()), tt.args.path)
			}
		})
	}
}

func TestAddNode(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name: "Success",
			path: "addnode",
		},
		{
			name:    "NodeExists",
			path:    "addnode",
			wantErr: true,
		},
		{
			name:    "ParentNotExists",
			path:    "othernode subnode",
			wantErr: true,
		},
		{
			name:    "Empty",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddNode(tt.path, tt.path)
			if tt.wantErr {
				if err == nil {
					t.Errorf("AddCommand() error = %v, wantErr %v", err, tt.wantErr)

				}
				assert.Nil(t, got)
				return
			} else if got == nil {
				assert.NotNil(t, got)
				return
			}
		})
	}
}

func TestExit(t *testing.T) {
	const exitCode = 10
	if os.Getenv("BE_EXITER") == "1" {
		SetExitCode(exitCode)
		Exit()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExit")
	cmd.Env = append(os.Environ(), "BE_EXITER=1")
	err := cmd.Run()

	exitError, ok := err.(*exec.ExitError)
	if ok && exitError.ProcessState.ExitCode() == exitCode {
		return
	}
	t.Fatalf("process ran with err %v, want exit status %d", err, exitCode)
}

func TestFatal(t *testing.T) {
	const exitCode = 1
	if os.Getenv("BE_EXITER") == "1" {
		Fatal(errors.New("some error"))
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "BE_EXITER=1")
	err := cmd.Run()

	exitError, ok := err.(*exec.ExitError)
	if ok && exitError.ProcessState.ExitCode() == exitCode {
		return
	}
	t.Fatalf("process ran with err %v, want exit status %d", err, exitCode)
}

func TestFatal_NoError(t *testing.T) {
	Fatal(nil)
	assert.True(t, true)
}

func TestFindCommand(t *testing.T) {
	var cmdFunc CommandFunc = func(args []string) error {
		return nil
	}

	_, _ = AddCommand("find-command", "find-command", cmdFunc)
	_, _ = AddCommand("find-command another-command", "another-command", cmdFunc)

	tests := []struct {
		name    string
		path    []string
		wantNil bool
	}{
		{
			name:    "Root",
			path:    []string{},
			wantNil: false,
		},
		{
			name:    "Primary",
			path:    []string{"find-command"},
			wantNil: false,
		},
		{
			name:    "Secondary",
			path:    []string{"find-command", "another-command"},
			wantNil: false,
		},
		{
			name:    "MissingPrimary",
			path:    []string{"finds-command"},
			wantNil: true,
		},
		{
			name:    "MissingSecondary",
			path:    []string{"find-command", "anothers-command"},
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindCommand(tt.path...)
			gotNil := got == nil

			if gotNil != tt.wantNil {
				t.Errorf("FindCommand() = %v, wantNil = %v", got, tt.wantNil)
			}
		})
	}
}

func TestRootCmd(t *testing.T) {
	assert.Equal(t, rootCmd, RootCmd())
}

func TestSetExitCode(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "Success",
			code: 0,
		},
		{
			name: "Failure",
			code: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetExitCode(tt.code)
			assert.Equal(t, tt.code, exitCode)
		})
	}
}
