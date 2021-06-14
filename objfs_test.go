package objfs

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFS(t *testing.T) {
	fs, err := NewFS(
		WithEndPoint(""),
		WithAccessKeySecret(""),
		WithBucketName(""),
		WithAccessKeyID(""))
	assert.Nil(t, err)
	f, err := fs.Open("blog/indexerror.PNG")
	assert.Nil(t, err)
	fi, err := f.Stat()
	assert.Nil(t, err)
	log.Printf("name: %s", fi.Name())
	log.Printf("size: %d", fi.Size())

	b := make([]byte, 100)
	size, err := f.Read(b)
	assert.Nil(t, err)
	assert.Equal(t, 100, size)
	log.Printf("%s", b)
}

func TestReadDir(t *testing.T) {
	fsys, err := NewFS(
		WithEndPoint(""),
		WithAccessKeySecret(""),
		WithBucketName(""),
		WithAccessKeyID(""))
	assert.Nil(t, err)
	res, err := fsys.ReadDir("blog")
	assert.Nil(t, err)
	for _, f := range res {
		t.Logf("name: %s, isdir %v", f.Name(), f.IsDir())
	}
}

