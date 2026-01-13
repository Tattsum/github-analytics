package domain

import "testing"

func TestNewUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		login     string
		userName  string
		createdAt string
		want      *User
	}{
		{
			name:      "正常なユーザー作成",
			login:     "testuser",
			userName:  "Test User",
			createdAt: "2020-01-01T00:00:00Z",
			want: &User{
				Login:     "testuser",
				Name:      "Test User",
				CreatedAt: "2020-01-01T00:00:00Z",
			},
		},
		{
			name:      "空の名前でも作成可能",
			login:     "testuser",
			userName:  "",
			createdAt: "2020-01-01T00:00:00Z",
			want: &User{
				Login:     "testuser",
				Name:      "",
				CreatedAt: "2020-01-01T00:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewUser(tt.login, tt.userName, tt.createdAt)
			if got.Login != tt.want.Login {
				t.Errorf("NewUser().Login = %v, want %v", got.Login, tt.want.Login)
			}

			if got.Name != tt.want.Name {
				t.Errorf("NewUser().Name = %v, want %v", got.Name, tt.want.Name)
			}

			if got.CreatedAt != tt.want.CreatedAt {
				t.Errorf("NewUser().CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
			}
		})
	}
}

func TestUser_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		user *User
		want bool
	}{
		{
			name: "有効なユーザー",
			user: NewUser("testuser", "Test User", "2020-01-01T00:00:00Z"),
			want: true,
		},
		{
			name: "ログイン名が空の場合は無効",
			user: NewUser("", "Test User", "2020-01-01T00:00:00Z"),
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.user.Login != ""
			if got != tt.want {
				t.Errorf("User.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
