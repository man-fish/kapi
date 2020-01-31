package models

type Group struct {
	ID					int64	`json:"id" sql:"id"`
	GroupName			string	`json:"group_name" sql:"group_name"`
	GroupDesc			string	`json:"group_desc" sql:"group_desc"`
	Type 				string	`json:"type" sql:"type"`
	CustomField			string	`json:"custom_field" sql:"custom_field"`
	CustomFieldEnabled string	`json:"custom_field_enable" sql:"custom_field_enable"`
	AddTime				string	`json:"add_time" sql:"add_time"`
	UpTime				string	`json:"up_time" sql:"up_time"`
}