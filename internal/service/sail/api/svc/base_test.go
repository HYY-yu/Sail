package svc

import (
	"context"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"reflect"
	"testing"
)

func TestBaseSvc_CheckStaffGroup(t *testing.T) {
	type args struct {
		ctx            context.Context
		projectGroupId int
	}
	tests := []struct {
		name  string
		args  args
		want  []int
		want1 model.Role
	}{
		{
			name: "Admin",
			args: args{
				ctx: context.WithValue(context.Background(), model.StaffGroupRelKey, []model.StaffGroup{
					{
						ProjectGroupID: 1,
						Role:           model.RoleAdmin,
					},
				}),
				projectGroupId: 1,
			},
			want:  nil,
			want1: model.RoleAdmin,
		},
		{
			name: "Member",
			args: args{
				ctx: context.WithValue(context.Background(), model.StaffGroupRelKey, []model.StaffGroup{
					{
						ProjectGroupID: 3,
						Role:           model.RoleOwner,
					},
					{
						ProjectGroupID: 5,
						Role:           model.RoleManager,
					},
				}),
				projectGroupId: 3,
			},
			want:  []int{3, 5},
			want1: model.RoleOwner,
		},
		{
			name: "NoHas",
			args: args{
				ctx: context.WithValue(context.Background(), model.StaffGroupRelKey, []int{
					1, 2, 3}),
				projectGroupId: 3,
			},
			want:  nil,
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BaseSvc{}
			got, got1 := s.CheckStaffGroup(tt.args.ctx, tt.args.projectGroupId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckStaffGroup() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckStaffGroup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
