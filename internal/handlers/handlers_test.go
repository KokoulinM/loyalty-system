package handlers

import "testing"

func TestHandlers_Register(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name  string
		query string
		body  string
		want  want
	}{
		{
			name:  "пользователь успешно аутентифицирован",
			query: "/api/user/register",
			body:  `{"first_name": "first_name", "last_name": "last_name", "login": "login", "password": "12345"}`,
			want: want{
				code:        200,
				contentType: "application/json; charset=utf-8",
			},
		},
		//{
		//	name:  "неверный формат запроса",
		//	query: "/api/user/register",
		//},
		//{
		//	name:  "неверная пара логин/пароль",
		//	query: "/api/user/register",
		//},
		//{
		//	name:  "внутренняя ошибка сервера",
		//	query: "/api/user/register",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
