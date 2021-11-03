package controller

func responseWithStatus(status int, msg string) UniversalResp {
	return UniversalResp{
		Status:  status,
		Message: msg,
	}
}

func responseWithData(data interface{}, count int, size int, offset int, msg string) UniversalRespByData {
	code := 0
	if msg != "" {
		code = 1
	}
	return UniversalRespByData{
		UniversalResp: UniversalResp{
			Status:  code,
			Message: msg,
		},
		Count:  count,
		Size:   size,
		Offset: offset,
		Data:   data,
	}
}
