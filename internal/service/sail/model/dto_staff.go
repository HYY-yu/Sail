package model

type StaffList struct {
	StaffID    int    `json:"staff_id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`

	Roles []StaffRole `json:"roles"`
}

type Role int

const (
	RoleAdmin Role = iota + 1
	RoleOwner
	RoleManager
)

func (r Role) String() string {
	return [...]string{"Unknown", "Admin", "Owner", "Manager"}[r]
}

type StaffRole struct {
	StaffGroupRelID  int    `json:"staff_group_rel_id"`
	ProjectGroupID   int    `json:"project_group_id"`
	ProjectGroupName string `json:"project_group_name"`
	Role             Role   `json:"role"`
	RoleInfo         string `json:"role_info"`
}

type AddStaff struct {
	Name string `json:"name" v:"required|regex:^[a-zA-Z][a-zA-Z1-9]{1,9}"` // 员工标识
}

type EditStaff struct {
	StaffID int     `json:"staff_id" v:"required"`
	Name    *string `json:"name" v:"regex:^[a-zA-Z][a-zA-Z1-9]{1,9}"`
}

type GrantStaff struct {
	StaffID        int  `json:"staff_id" v:"required"`
	Role           Role `json:"role" v:"required"`
	ProjectGroupID int  `json:"project_group_id"`
}

type StaffGroup struct {
	ProjectGroupID int  `json:"project_group_id"`
	Role           Role `json:"role"`
}
