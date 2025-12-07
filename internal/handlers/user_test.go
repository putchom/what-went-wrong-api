package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"what-went-wrong-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetUsers_Success(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	testUsers := []models.User{
		{ID: uuid.New(), Auth0ID: "auth0|1", Name: "山田太郎", Email: "yamada@example.com"},
		{ID: uuid.New(), Auth0ID: "auth0|2", Name: "佐藤花子", Email: "sato@example.com"},
		{ID: uuid.New(), Auth0ID: "auth0|3", Name: "鈴木一郎", Email: "suzuki@example.com"},
	}

	for _, user := range testUsers {
		result := db.Create(&user)
		assert.NoError(t, result.Error, "テストデータの挿入に失敗しました")
	}

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	c.Request = req

	handler := GetUsers(db)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code, "ステータスコードが200であるべきです")

	var responseUsers []models.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	assert.NoError(t, err, "レスポンスのパースに失敗しました")

	assert.Equal(t, 3, len(responseUsers), "3件のユーザーが返されるべきです")
	// Order is not guaranteed, so we should check existence or sort.
	// But for now let's assume insertion order or check one.
	// Actually, with Postgres, order is not guaranteed without Order by.
	// But let's check if we can find "山田太郎".
	found := false
	for _, u := range responseUsers {
		if u.Name == "山田太郎" {
			found = true
			assert.Equal(t, "yamada@example.com", u.Email)
			break
		}
	}
	assert.True(t, found, "山田太郎が見つかりませんでした")
}

func TestGetUsers_EmptyDatabase(t *testing.T) {
	db, cleanup := SetupTestDB(t)
	defer cleanup()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	c.Request = req

	handler := GetUsers(db)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code, "ステータスコードが200であるべきです")

	var responseUsers []models.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	assert.NoError(t, err, "レスポンスのパースに失敗しました")

	assert.Equal(t, 0, len(responseUsers), "空の配列が返されるべきです")
}

func TestGetUsers_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		setupDB        func(*gorm.DB)
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "正常系: 複数のユーザーが存在する",
			setupDB: func(db *gorm.DB) {
				db.Create(&models.User{ID: uuid.New(), Auth0ID: "auth0|1", Name: "User1", Email: "user1@example.com"})
				db.Create(&models.User{ID: uuid.New(), Auth0ID: "auth0|2", Name: "User2", Email: "user2@example.com"})
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "正常系: ユーザーが存在しない",
			setupDB:        func(db *gorm.DB) {},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "正常系: 1件のユーザーが存在する",
			setupDB: func(db *gorm.DB) {
				db.Create(&models.User{ID: uuid.New(), Auth0ID: "auth0|3", Name: "SingleUser", Email: "single@example.com"})
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := SetupTestDB(t)
			defer cleanup()

			if tt.setupDB != nil {
				tt.setupDB(db)
			}
			gin.SetMode(gin.TestMode)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/api/v1/users", nil)
			c.Request = req

			handler := GetUsers(db)
			handler(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseUsers []models.User
			err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCount, len(responseUsers))
		})
	}
}
