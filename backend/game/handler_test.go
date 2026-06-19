package game

import "testing"

func Test_filterText(t *testing.T) {
	type args struct {
		original string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			"Пустой ответ",
			args{"TESTовый, -"},
			"",
			false,
		},
		{
			"Ответ без имени бота",
			args{"правильный ответ"},
			"правильный ответ",
			true,
		},
		{
			"Ответ капсом",
			args{"правильный ОТВЕТ"},
			"правильный ответ",
			true,
		},
		{
			"Ответ c именем бота",
			args{"TESTовый, мой ответ"},
			"мой ответ",
			true,
		},
		{
			"Ответ с запятыми",
			args{"TESTовый, розовый, бирюзовый"},
			"розовый, бирюзовый",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := filterText(tt.args.original)
			if got != tt.want {
				t.Errorf("filterText() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("filterText() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
