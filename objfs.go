package objfs

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type FS struct {
	client *oss.Client
	bucket *oss.Bucket
}

func NewFS(opt ...Opt) (FS, error) {
	var (
		client *oss.Client
		bucket *oss.Bucket
		err    error
	)
	var newFsOpts newFsOptions
	for _, f := range opt {
		err = f(&newFsOpts)
		if err != nil {
			return FS{}, err
		}
	}
	if err = newFsOpts.Validate(); err != nil {
		return FS{}, err
	}

	if client, err = oss.New(
		newFsOpts.endpoint,
		newFsOpts.accessKeyID,
		newFsOpts.accessKeySecret,
	); err != nil {
		return FS{}, err
	}
	if bucket, err = client.Bucket(newFsOpts.bucketName); err != nil {
		return FS{}, fmt.Errorf("failed to find bucket %s, err: %w", newFsOpts.bucketName, err)
	}
	return FS{
		client,
		bucket,
	}, nil
}

func (fileSys FS) ReadDir(name string) ([]fs.DirEntry, error) {
	//append backslash if needed
	if name == "" {
		name = "/"
	}
	if !strings.HasSuffix(name, string(os.PathSeparator)) {
		name = fmt.Sprintf("%s%c", name, os.PathSeparator)
	}
	//TODO: check if res IsTruncated
	var opt []oss.Option
	opt = append(opt, oss.Delimiter(string(filepath.Separator)))
	if name != "/" {
		opt = append(opt, oss.Prefix(name))
	}
	listRes, err := fileSys.bucket.ListObjectsV2(opt...)
	if err != nil {
		return nil, err
	}
	var res []fs.DirEntry
	for _, d := range listRes.CommonPrefixes {
		if strings.EqualFold(d, name) {
			continue
		}
		res = append(res, dEntry{
			prefix: name,
			name:   d,
		})
	}
	for _, o := range listRes.Objects {
		if strings.EqualFold(o.Key, name) {
			continue
		}
		res = append(res, dEntry{
			prefix: name,
			name:   o.Key,
		})
	}
	return res, nil
}

type dEntry struct {
	prefix string
	name   string
}

func (d dEntry) Name() string {
	return strings.TrimSuffix(strings.TrimLeft(d.name, d.prefix), string(os.PathSeparator))
}

func (d dEntry) IsDir() bool {
	return strings.HasSuffix(d.name, string(os.PathSeparator))
}

func (d dEntry) Type() fs.FileMode {
	panic("implement me")
}

func (d dEntry) Info() (fs.FileInfo, error) {
	panic("implement me")
}

type ObjectFile struct {
	fs   FS
	name string
	r    io.ReadCloser
}

type ObjectFileStat struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
}

func (o ObjectFileStat) Name() string {
	return o.name
}

func (o ObjectFileStat) Size() int64 {
	return o.size
}

func (o ObjectFileStat) Mode() fs.FileMode {
	return o.mode
}

func (o ObjectFileStat) ModTime() time.Time {
	return o.modTime
}

func (o ObjectFileStat) IsDir() bool {
	return false
}

func (o ObjectFileStat) Sys() interface{} {
	return nil
}

func (file *ObjectFile) Stat() (fs.FileInfo, error) {
	meta, err := file.fs.bucket.GetObjectDetailedMeta(file.name)
	if err != nil {
		return nil, err
	}
	sizeStr := meta.Get("Content-Length")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return ObjectFileStat{}, err
	}
	modTimeStr := meta.Get("Last-Modified")
	lastModified, err := time.Parse(http.TimeFormat, modTimeStr)
	if err != nil {
		return ObjectFileStat{}, err
	}
	stat := ObjectFileStat{
		name:    file.name,
		size:    int64(size),
		mode:    0,
		modTime: lastModified,
	}
	return stat, nil
}

func (file *ObjectFile) Read(bytes []byte) (int, error) {
	return file.r.Read(bytes)
}

func (file *ObjectFile) Close() error {
	return file.r.Close()
}

// Open a/b/c
func (fileSys FS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	r, err := fileSys.bucket.GetObject(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return &ObjectFile{
		fs:   fileSys,
		name: name,
		r:    r,
	}, nil
}
