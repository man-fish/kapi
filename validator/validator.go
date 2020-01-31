package validators

type LoginValidator struct {
	Email    			string 	`kapi:"email" validate:"nonzero"`
	Password 			string 	`kapi:"password" validate:"min=6,max=30"`
}


type RegisterValidator struct {
	LoginValidator
	Username 			string	`kapi:"username" validate:"min=3,max=30"`
}

type CreateGroupValidator struct {
	GroupName			string	`kapi:"group_name" validate:"min=1,max=40"`
	GroupDesc			string	`kapi:"group_desc" validate:"min=1,max=200"`
	Type 				string	`kapi:"type" validate:"min=1,max=10"`
}

type UpdateGroupValidator struct {
	ID					int64	`kapi:"id"`
	GroupName			string	`kapi:"group_name" validate:"min=1,max=40"`
	GroupDesc			string	`kapi:"group_desc" validate:"min=1,max=200"`
	CustomFieldEnabled 	string	`kapi:"custom_field_enable" validate:"min=1,max=1"`
	CustomField			string	`kapi:"custom_field" validate:"min=1,max=50"`
}

type AddMemberValidator struct {
	Gid					int64   `kapi:"gid" validate:"min=1"`
	Role 				string	`kapi:"role" validate:"min=1,max=10"`
	Email 				string	`kapi:"email" validate:"nonzero"`
}

type DelMemberValidator struct {
	Uid					int64   `kapi:"uid" validate:"min=1"`
	Gid 				int64	`kapi:"gid" validate:"min=1"`
}

type UpdateMemberValidator struct {
	DelMemberValidator
	ID					int64 `kapi:"id" validate:"min=1"`
	Role 				string `kapi:"role" validate:"min=1,max=10"`
}