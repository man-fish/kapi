package utils

func MakeDefaultRes(status int64, msg interface{}, data interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	res["status"] =	status
	res["msg"] = msg
	res["data"] = data
	return res
}
