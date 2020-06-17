package cla_groups

import "testing"

func Test_validateEnrollProject(t *testing.T) {
	type args struct {
		cgmap           map[string]int
		claGroupID      string
		projectSFIDList []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// new cla group when no other cla group is present
		{
			name: "new foundation level cla-group",
			args: args{
				cgmap:           map[string]int{},
				claGroupID:      "",
				projectSFIDList: []string{"1", "2"},
			},
		},
		{
			name: "new project level cla-group",
			args: args{
				cgmap:           map[string]int{},
				claGroupID:      "",
				projectSFIDList: []string{"1"},
			},
		},
		{
			name: "new project level cla-group with no projects",
			args: args{
				cgmap:           map[string]int{},
				claGroupID:      "",
				projectSFIDList: []string{},
			},
		},
		// new cla group when other cla group is present
		{
			name: "new foundation level cla-group, when other foundation level cla-group is present",
			args: args{
				cgmap:           map[string]int{"cg1": 2},
				claGroupID:      "",
				projectSFIDList: []string{"1", "2"},
			},
			wantErr: true,
		},
		{
			name: "new foundation level cla-group, when other project level cla-group is present",
			args: args{
				cgmap:           map[string]int{"cg1": 1},
				claGroupID:      "",
				projectSFIDList: []string{"1", "2"},
			},
			wantErr: true,
		},
		{
			name: "new project level cla-group, when other foundation level cla-group is present",
			args: args{
				cgmap:           map[string]int{"cg1": 2},
				claGroupID:      "",
				projectSFIDList: []string{"1"},
			},
			wantErr: true,
		},
		{
			name: "new project level cla-group, when other project level cla-group is present",
			args: args{
				cgmap:           map[string]int{"cg1": 1},
				claGroupID:      "",
				projectSFIDList: []string{"1"},
			},
			wantErr: false,
		},
		// add projects to existing cla groups
		{
			name: "enroll projects in existing foundation level cla_group",
			args: args{
				cgmap:           map[string]int{"cg1": 2},
				claGroupID:      "cg1",
				projectSFIDList: []string{"1", "2"},
			},
			wantErr: false,
		},
		{
			name: "enroll projects in existing project level cla_group, when no other project level cla group present",
			args: args{
				cgmap:           map[string]int{"cg1": 1},
				claGroupID:      "cg1",
				projectSFIDList: []string{"1", "2"},
			},
			wantErr: false,
		},
		{
			name: "enroll projects in existing project level cla_group, when other project level cla group present",
			args: args{
				cgmap:           map[string]int{"cg1": 1, "cg2": 1},
				claGroupID:      "cg1",
				projectSFIDList: []string{"1", "2"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEnrollProject(tt.args.cgmap, tt.args.claGroupID, tt.args.projectSFIDList); (err != nil) != tt.wantErr {
				t.Errorf("validateEnrollProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
