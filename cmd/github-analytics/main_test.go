package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetUsersFromStrings は文字列からユーザーリストを取得する関数のテストです.
func TestGetUsersFromStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		usersStr string
		want     []string
	}{
		{
			name:     "単一ユーザー",
			usersStr: "user1",
			want:     []string{"user1"},
		},
		{
			name:     "複数ユーザー（カンマ区切り）",
			usersStr: "user1,user2,user3",
			want:     []string{"user1", "user2", "user3"},
		},
		{
			name:     "空白を含むユーザー名",
			usersStr: "user1, user2 , user3",
			want:     []string{"user1", "user2", "user3"},
		},
		{
			name:     "空文字列は除外",
			usersStr: "user1,,user2",
			want:     []string{"user1", "user2"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := getUsersFromStrings(tt.usersStr)
			assert.Equal(t, tt.want, got, "getUsersFromStrings() should return correct users")
		})
	}
}

// TestMergeUsers はユーザーリストをマージする関数のテストです.
func TestMergeUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		orgMembers     []string
		specifiedUsers []string
		want           []string
		wantAdded      int
	}{
		{
			name:           "組織メンバーのみ",
			orgMembers:     []string{"user1", "user2", "user3"},
			specifiedUsers: []string{},
			want:           []string{"user1", "user2", "user3"},
			wantAdded:      0,
		},
		{
			name:           "指定ユーザーのみ",
			orgMembers:     []string{},
			specifiedUsers: []string{"user1", "user2"},
			want:           []string{"user1", "user2"},
			wantAdded:      2,
		},
		{
			name:           "組織メンバーと指定ユーザー（重複なし）",
			orgMembers:     []string{"user1", "user2"},
			specifiedUsers: []string{"user3", "user4"},
			want:           []string{"user1", "user2", "user3", "user4"},
			wantAdded:      2,
		},
		{
			name:           "組織メンバーと指定ユーザー（重複あり）",
			orgMembers:     []string{"user1", "user2"},
			specifiedUsers: []string{"user2", "user3"},
			want:           []string{"user1", "user2", "user3"},
			wantAdded:      1,
		},
		{
			name:           "すべて重複",
			orgMembers:     []string{"user1", "user2"},
			specifiedUsers: []string{"user1", "user2"},
			want:           []string{"user1", "user2"},
			wantAdded:      0,
		},
		{
			name:           "空のリスト",
			orgMembers:     []string{},
			specifiedUsers: []string{},
			want:           []string{},
			wantAdded:      0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotAdded := mergeUsers(tt.orgMembers, tt.specifiedUsers)
			assert.Equal(t, tt.want, got, "mergeUsers() should return correct merged users")
			assert.Equal(t, tt.wantAdded, gotAdded, "mergeUsers() should return correct added count")
		})
	}
}
