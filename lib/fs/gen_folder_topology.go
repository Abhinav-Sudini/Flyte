package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func formatSize(b int64) string {
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	}
	if b < 1024*1024 {
		return fmt.Sprintf("%d KB", b/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
}

func formatDate(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return "Today"
	}
	if t.Year() == now.Year() && t.YearDay()+1 == now.YearDay() {
		return "Yesterday"
	}
	return t.Format("Jan 02, 2006")
}

func typeOf(name string, isDir bool) string {
	if isDir {
		return "folder"
	}
	switch strings.ToLower(filepath.Ext(name)) {
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "doc"
	case ".xls", ".xlsx":
		return "sheet"
	case ".png", ".jpg", ".jpeg":
		return "image"
	default:
		return strings.TrimPrefix(filepath.Ext(name), ".")
	}
}
func formatPar(par int) string {
	if par == -1 {
		return "null"
	} else {
		return fmt.Sprintf("%d", par)
	}
}

type Entry struct {
	Name    string
	Size    int64
	ModTime time.Time
	isDir   bool
}

func get_all_files_ase(path string) ([]Entry, error) {
	var list []Entry

	entries, err := os.ReadDir(path)
	if err != nil {
		return list, err
	}

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		list = append(list, Entry{
			Name:    e.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			isDir:   info.IsDir(),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ModTime.After(list[j].ModTime) // descending
	})
	return list, nil
}
func dfs_add_files(root string, par int, out []string) ([]string, error) {
	entries, err := get_all_files_ase(root)
	if err != nil {
		return out, err
	}
	for _, entr := range entries {
		tem := fmt.Sprintf("{ \"id\": %d, \"parentId\": %s, \"name\": \"%s\", \"type\": \"%s\", \"size\": \"%s\", \"date\": \"%s\"}", len(out)+1, formatPar(par), entr.Name, typeOf(entr.Name, entr.isDir), formatSize(entr.Size), formatDate(entr.ModTime))
		out = append(out, tem)
		if entr.isDir {
			// fmt.Println(entr.Name)
			out, _ = dfs_add_files(root+"/"+entr.Name, len(out), out)
		}
	}
	return out, nil
}
func DirToJSON(root string) (string, error) {
	var out []string
	out, err := dfs_add_files(root, -1, out)
	jsonStr := "[\n  " + strings.Join(out, ",\n  ") + "\n]"
	//fmt.Println(jsonStr)
	return jsonStr, err
}

// func main(){
// 	out,err := DirToJSON("/home/abhi/code/projects/Flyte")
// 	fmt.Println(out,err)
// }
