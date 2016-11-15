package scripts

import (
    "time"
    "github.com/lxn/walk"
    "path/filepath"
    "os"
    "strings"
    "fmt"
    "log"
)

type Directory struct {
    name     string
    parent   *Directory
    children []*Directory
}

func NewDirectory(name string, parent *Directory) *Directory {
    return &Directory{name: name, parent: parent}
}

var _ walk.TreeItem = new(Directory)

func (d *Directory) Text() string {
    return d.name
}

func (d *Directory) Parent() walk.TreeItem {
    if d.parent == nil {
        // We can't simply return d.parent in this case, because the interface
        // value then would not be nil.
        return nil
    }
    return d.parent
}

func (d *Directory) ChildCount() int {
    if d.children == nil {
        // It seems this is the first time our child count is checked, so we
        // use the opportunity to populate our direct children.
        if err := d.ResetChildren(); err != nil {
            log.Print(err)
        }
    }
    return len(d.children)
}

func (d *Directory) ChildAt(index int) walk.TreeItem {
    return d.children[index]
}

func (d *Directory) Image() interface{} {
    return d.Path()
}

func (d *Directory) ResetChildren() error {
    d.children = nil

    dirPath := d.Path()

    if err := filepath.Walk(d.Path(), func(path string, info os.FileInfo, err error) error {
        if err != nil {
            if info == nil {
                return filepath.SkipDir
            }
        }

        name := info.Name()

        if !info.IsDir() || path == dirPath || shouldExclude(name) {
            return nil
        }

        d.children = append(d.children, NewDirectory(name, d))

        return filepath.SkipDir
    }); err != nil {
        return err
    }

    return nil
}

func (d *Directory) Path() string {
    elems := []string{d.name}

    dir, _ := d.Parent().(*Directory)

    for dir != nil {
        elems = append([]string{dir.name}, elems...)
        dir, _ = dir.Parent().(*Directory)
    }

    return filepath.Join(elems...)
}

type DirectoryTreeModel struct {
    walk.TreeModelBase
    roots []*Directory
}

var _ walk.TreeModel = new(DirectoryTreeModel)

func NewDirectoryTreeModel() (*DirectoryTreeModel, error) {
    model := new(DirectoryTreeModel)
    model.roots = nil
    //drives, err := walk.DriveNames()
    //if err != nil {
    //    return nil, err
    //}
    path := NewDirectory("F:\\", nil)
    model.roots = append(model.roots, NewDirectory("Learn", path))

    //for _, drive := range drives {
    //    switch drive {
    //    case "A:\\", "B:\\":
    //        continue
    //    }
    //
    //    model.roots = append(model.roots, NewDirectory(drive, nil))
    //}

    return model, nil
}

func (*DirectoryTreeModel) LazyPopulation() bool {
    // We don't want to eagerly populate our tree view with the whole file system.
    return true
}

func (m *DirectoryTreeModel) RootCount() int {
    return len(m.roots)
}

func (m *DirectoryTreeModel) RootAt(index int) walk.TreeItem {
    return m.roots[index]
}

type FileInfo struct {
    Name     string
    Path     string
    Size     string
    Modified time.Time

}

type FileInfoModel struct {
    walk.SortedReflectTableModelBase
    dirPath string
    items   []*FileInfo
}

var _ walk.ReflectTableModel = new(FileInfoModel)

func NewFileInfoModel() *FileInfoModel {
    return new(FileInfoModel)
}

func (m *FileInfoModel) Items() interface{} {
    return m.items
}

func (m *FileInfoModel) SetDirPath(dirPath string) error {
    m.dirPath = dirPath
    m.items = nil

    if err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            if info == nil {
                return filepath.SkipDir
            }
        }

        name := info.Name()

        if path == dirPath || shouldExclude(name) {
            return nil
        }
        size := ""
        if  int(info.Size() / (1024 * 1024)) > 1  {
            size = fmt.Sprintf( "%d MB", int(info.Size() / (1024 * 1024)) )
        } else {
            size = fmt.Sprintf( "%d KB", int(info.Size() / 1024) )
        }
        item := &FileInfo{
            Name:     name,
            Path:     "",
            Size:     size,
            Modified: info.ModTime(),
        }

        m.items = append(m.items, item)

        if info.IsDir() {
            return filepath.SkipDir
        }

        return nil
    }); err != nil {
        return err
    }

    m.PublishRowsReset()

    return nil
}

func (m *FileInfoModel) Image(row int) interface{} {
    return filepath.Join(m.dirPath, m.items[row].Name)
}

func shouldExclude(name string) bool {
    switch name {
    case "System Volume Information", "pagefile.sys", "swapfile.sys":
        return true
    }

    return false
}


//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func (m *FileInfoModel) WalkDir(dirPth string, suffixArr []string) (files []string, names []string, err error) {
    //m.dirPath = dirPth
    m.items = nil

    files = make([]string, 0, 100)                                                       //文件路径
    names = make([]string, 0, 100)
    //文件名称

    err = filepath.Walk(dirPth, func(filename string, info os.FileInfo, err error) error { //遍历目录
        if err != nil { //忽略错误
            return err
        }
        if info.IsDir() { // 忽略目录
            return nil
        }
        //忽略后缀匹配的大小写
        for _, v := range suffixArr {
            suffix := "." + strings.ToUpper(v) //忽略后缀匹配的大小写
            if strings.HasSuffix(strings.ToUpper(info.Name()), suffix) {
                size := ""
                if  int(info.Size() / (1024 * 1024)) > 1  {
                    size = fmt.Sprintf( "%d MB", int(info.Size() / (1024 * 1024)) )
                } else {
                    size = fmt.Sprintf( "%d KB", int(info.Size() / 1024) )
                }
                item := &FileInfo{
                    Name:     info.Name(),
                    Path:     filepath.Dir(filename),
                    Size:     size,
                    Modified: info.ModTime(),
                }
                m.items = append(m.items, item)
            }
        }
        m.PublishRowsReset()
        return nil
    })
    return files, names, err
}
