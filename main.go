// ファイルの変更を確認するパッケージ
package checkmodfile

import (
	"os"
	"time"
)

type File struct {
	Name    string
	ModTime time.Time
}

// 管理対象に登録
func RegistFile(filename string) (*File, error) {
	f := new(File)
	f.Name = filename
	err := f.UpdateMod()
	if err != nil {
		return nil, err
	}
	return f, nil
}

// ModTimeを新しい物に更新
func (f *File) UpdateMod() error {
	fInfo, err := os.Stat(f.Name)
	if err != nil {
		return err
	}
	f.ModTime = fInfo.ModTime()
	return nil
}

//
func (f *File) IsLatest() (bool, error) {
	fInfo, err := os.Stat(f.Name)
	if err != nil {
		return false, err
	}
	return f.ModTime == fInfo.ModTime(), nil
}
