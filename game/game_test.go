package game

import "testing"

func Test_isCorrectAnswer(t *testing.T) {
	type args struct {
		gameAnswers []string
		userAnswer  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Ответ одно число",
			args{[]string{"1"}, "1"},
			true,
		},
		{
			"Ответ одно слово кириллицей",
			args{[]string{"ответ"}, "ответ"},
			true,
		},
		{
			"Ответ несколько слов кириллицей",
			args{[]string{"правильный ответ"}, "правильный ответ"},
			true,
		},
		{
			"Ответ несколько слов кириллицей с запятыми",
			args{[]string{"правильный, ответ"}, "правильный, ответ"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCorrectAnswer(tt.args.gameAnswers, tt.args.userAnswer); got != tt.want {
				t.Errorf("isCorrectAnswer() = %v, want %v", got, tt.want)
			}
		})
	}
}
