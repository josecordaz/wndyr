package main

type Error struct {
	Message string `json:"message"`
}

type Photo struct {
	ImgSrc string `json:"img_src"`
}

type Data struct {
	Photos []*Photo `json:"photos"`
	Err    *Error   `json:"error"`
}
