package store

import (
	"Flyte/util"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

const (
	MAX_SAVE_FILE_SIZE = 10 << 20
)

func SaveNewFile(user_id int, upload_dir string, new_file multipart.File, handler *multipart.FileHeader) error {

	server_upload_dir := getServerUploadDirPath(user_id, upload_dir)

	if util.Dir_not_exist(server_upload_dir) {
		return errors.New("upload dir does not exist")
	}

	if handler.Size > MAX_SAVE_FILE_SIZE {
		return errors.New("to big of a file")
	}

	fileLocation := filepath.Join(server_upload_dir, handler.Filename)
	log.Println("store : ",fileLocation)
	dst, err := os.Create(fileLocation)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, new_file); err != nil {
		return err
	}
	return nil
}

func SaveNewFolder(user_id int, upload_path string) error {
	server_upload_dir := getServerUploadDirPath(user_id, upload_path)
	if util.Dir_not_exist(server_upload_dir) {
		return errors.New("dir does not exist")
	}
	log.Println("store new folde :",server_upload_dir)
	return os.Mkdir(filepath.Join(server_upload_dir,"NewFolder"), 0755)
}

func RenameFile(user_id int, upload_dir string, old_file_name string, new_file_name string) error {
	server_upload_dir := getServerUploadDirPath(user_id, upload_dir)
	if util.Dir_not_exist(server_upload_dir) {
		return errors.New("dir does not exist")
	}
	if util.File_not_exist(filepath.Join(server_upload_dir, old_file_name)) &&
		 util.Dir_not_exist(filepath.Join(server_upload_dir, old_file_name)) {
		return errors.New("file/Dir does not exist")
	}
	return os.Rename(filepath.Join(server_upload_dir, old_file_name), filepath.Join(server_upload_dir, new_file_name))
}

func DeleteFile(user_id int, upload_dir string, file_name string) error {
	server_upload_dir := getServerUploadDirPath(user_id, upload_dir)
	server_file_path := filepath.Join(server_upload_dir, file_name)
	if util.File_not_exist(server_file_path) && util.Dir_not_exist(server_file_path) {
		return errors.New("file not exist")
	}
	return os.RemoveAll(server_file_path)
}

func getServerUploadDirPath(user_id int, upload_dir string) string {
	USER_SERVER_FS_ROOT := filepath.Join("./ALL_FILES/", util.Int_to_string(user_id))
	return filepath.Join(USER_SERVER_FS_ROOT, upload_dir)
}
