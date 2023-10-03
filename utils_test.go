package main

import (
	"reflect"
	"testing"
)

func Test_setFromConfig(t *testing.T) {
	type args struct {
		conf *NotifyConfig
		box  Box
	}
	tests := []struct {
		name string
		args args
		want Box
	}{
		{
			name: "success - set actions from configuration",
			args: args{
				conf: &NotifyConfig{
					Alias:         "test",
					OnNewMail:     "echo new",
					OnNewMailPost: "echo new-post",
				},
				box: Box{
					Alias: "test",
				},
			},
			want: Box{
				Alias:         "test",
				OnNewMail:     "echo new",
				OnNewMailPost: "echo new-post",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setFromConfig(tt.args.conf, tt.args.box); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setFromConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_retrieveCmd(t *testing.T) {
	type args struct {
		conf *NotifyConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *NotifyConfig
		wantErr bool
	}{
		{
			name: "success - retrieve password from CMD",
			args: args{
				conf: &NotifyConfig{
					PasswordCMD: "echo password",
				},
			},
			want: &NotifyConfig{
				Password:    "password",
				PasswordCMD: "echo password",
			},
		},
		{
			name: "success - retrieve username from CMD",
			args: args{
				conf: &NotifyConfig{
					UsernameCMD: "echo user@example.com",
				},
			},
			want: &NotifyConfig{
				Username:    "user@example.com",
				UsernameCMD: "echo user@example.com",
			},
		},
		{
			name: "success - retrieve hostname from CMD",
			args: args{
				conf: &NotifyConfig{
					HostCMD: "echo localhost",
				},
			},
			want: &NotifyConfig{
				Host:    "localhost",
				HostCMD: "echo localhost",
			},
		},
		{
			name: "failure - retrieve password from CMD",
			args: args{
				conf: &NotifyConfig{
					PasswordCMD: "exit 1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failure - retrieve username from CMD",
			args: args{
				conf: &NotifyConfig{
					UsernameCMD: "exit 1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failure - retrieve hostname from CMD",
			args: args{
				conf: &NotifyConfig{
					HostCMD: "exit 1",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := retrieveCmd(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("retrieveCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
