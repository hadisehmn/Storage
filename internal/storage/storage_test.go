package storage

import (
	"context"
	"go-practice/STORAGE/internal/auth"
	"go-practice/STORAGE/internal/auth/user"
	"go-practice/STORAGE/internal/database"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	models "go-practice/STORAGE/internal/model"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(
		filepath.Join(root, "../../server/.env"),
	)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestUpload(t *testing.T) {

	file, err := os.Open("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	header := &multipart.FileHeader{
		Filename: "test.txt",
	}

	db, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(context.Background())

	userRepository := user.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository)

	userID, err := authService.SignUp(models.SignUpRequest{
		Name:     "Upload Test User",
		Email:    "upload-test-" + uuid.New().String() + "@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatal(err)
	}

	repository := NewStorageRepository(db)
	service := NewStorageService(repository)

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "user id required",
			userID:  "",
			wantErr: true,
		},
		{
			name:    "upload successfully",
			userID:  userID,
			wantErr: false,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			err := service.Upload(test.userID, file, header)

			if test.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			files, err := repository.FindByUserID(test.userID)
			if err != nil {
				t.Fatalf("failed to find uploaded file: %v", err)
			}

			if len(files) != 1 {
				t.Fatalf("expected 1 file, got %d", len(files))
			}

			if files[0].UserID != test.userID {
				t.Errorf(
					"expected user id %q, got %q",
					test.userID,
					files[0].UserID,
				)
			}

			if files[0].FileName != "test.txt" {
				t.Errorf(
					"expected file name %q, got %q",
					"test.txt",
					files[0].FileName,
				)
			}

			if _, err := os.Stat(files[0].FilePath); err != nil {
				t.Errorf(
					"uploaded file does not exist: %v",
					err,
				)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(context.Background())

	userRepository := user.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository)

	userID, err := authService.SignUp(models.SignUpRequest{
		Name:     "Delete Test User",
		Email:    "delete-test-" + uuid.New().String() + "@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatal(err)
	}

	repository := NewStorageRepository(db)
	service := NewStorageService(repository)

	file, err := os.Open("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	header := &multipart.FileHeader{
		Filename: "test.txt",
	}

	err = service.Upload(userID, file, header)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	files, err := repository.FindByUserID(userID)
	if err != nil {
		t.Fatalf("failed to find file: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	uploadedFile := files[0]

	_, err = os.Stat(uploadedFile.FilePath)
	if err != nil {
		t.Fatalf("file does not exist: %v", err)
	}

	err = service.DeleteFile(uploadedFile.ID, userID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	files, err = repository.FindByUserID(userID)
	if err != nil {
		t.Fatalf("failed to find files after delete: %v", err)
	}

	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}

	_, err = os.Stat(uploadedFile.FilePath)

	if !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted, but it still exists")
	}
}

func TestDownloadFile(t *testing.T) {
	db, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(context.Background())

	userRepository := user.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository)

	userID, err := authService.SignUp(models.SignUpRequest{
		Name:     "Download Test User",
		Email:    "download-test-" + uuid.New().String() + "@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatal(err)
	}

	repository := NewStorageRepository(db)
	service := NewStorageService(repository)

	file, err := os.Open("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	header := &multipart.FileHeader{
		Filename: "test.txt",
	}

	err = service.Upload(userID, file, header)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	files, err := repository.FindByUserID(userID)
	if err != nil {
		t.Fatalf("failed to find file: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	uploadedFile := files[0]

	result, err := service.DownloadFile(
		uploadedFile.ID,
		userID,
	)

	if err != nil {
		t.Fatalf("download failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected file, got nil")
	}

	if result.ID != uploadedFile.ID {
		t.Errorf(
			"expected file ID %q, got %q",
			uploadedFile.ID,
			result.ID,
		)
	}

	if result.UserID != userID {
		t.Errorf(
			"expected user ID %q, got %q",
			userID,
			result.UserID,
		)
	}

	if result.FileName != "test.txt" {
		t.Errorf(
			"expected file name %q, got %q",
			"test.txt",
			result.FileName,
		)
	}

	if _, err := os.Stat(result.FilePath); err != nil {
		t.Errorf(
			"downloaded file does not exist: %v",
			err,
		)
	}
}

func TestGetUserFiles(t *testing.T) {
	db, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(context.Background())

	userRepository := user.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository)

	userID, err := authService.SignUp(models.SignUpRequest{
		Name:     "Get Files Test User",
		Email:    "get-files-test-" + uuid.New().String() + "@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatal(err)
	}

	repository := NewStorageRepository(db)
	service := NewStorageService(repository)

	file, err := os.Open("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	header := &multipart.FileHeader{
		Filename: "test.txt",
	}

	err = service.Upload(userID, file, header)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	files, err := service.GetUserFiles(userID)
	if err != nil {
		t.Fatalf("failed to get user files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	if files[0].UserID != userID {
		t.Errorf(
			"expected user ID %q, got %q",
			userID,
			files[0].UserID,
		)
	}

	if files[0].FileName != "test.txt" {
		t.Errorf(
			"expected file name %q, got %q",
			"test.txt",
			files[0].FileName,
		)
	}

	if _, err := os.Stat(files[0].FilePath); err != nil {
		t.Errorf(
			"file does not exist: %v",
			err,
		)
	}
}
