package domain

// User はGitHubユーザーを表すエンティティです.
type User struct {
	Login     string
	Name      string
	CreatedAt string // GitHubアカウント作成日時
}

// NewUser は新しいUserエンティティを作成します.
func NewUser(login, name, createdAt string) *User {
	return &User{
		Login:     login,
		Name:      name,
		CreatedAt: createdAt,
	}
}
