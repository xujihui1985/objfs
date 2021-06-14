package objfs

import (
	"github.com/xujihui1985/objfs/internal/xerror"
)

type newFsOptions struct {
	bucketName      string
	endpoint        string
	accessKeyID     string
	accessKeySecret string
}

type Opt func(opt *newFsOptions) error

func WithBucketName(bucketName string) Opt {
	return func(opt *newFsOptions) error {
		opt.bucketName = bucketName
		return nil
	}
}
func WithEndPoint(endpoint string) Opt {
	return func(opt *newFsOptions) error {
		opt.endpoint = endpoint
		return nil
	}
}
func WithAccessKeyID(accessKeyID string) Opt {
	return func(opt *newFsOptions) error {
		opt.accessKeyID = accessKeyID
		return nil
	}
}
func WithAccessKeySecret(accessKeySecret string) Opt {
	return func(opt *newFsOptions) error {
		opt.accessKeySecret = accessKeySecret
		return nil
	}
}

func (opt newFsOptions) Validate() error {
	e := new(xerror.ErrorCollector)
	if opt.endpoint == "" {
		e.Collect(xerror.NewInvalidArgsErr("endpoint", "endpoint is required"))
	}
	if opt.accessKeyID == "" {
		e.Collect(xerror.NewInvalidArgsErr("accessKeyID", "accessKeyID is required"))
	}
	if opt.accessKeySecret == "" {
		e.Collect(xerror.NewInvalidArgsErr("accessKeySecret", "accessKeySecret is required"))
	}
	if opt.bucketName == "" {
		e.Collect(xerror.NewInvalidArgsErr("bucketName", "bucketName is required"))
	}
	if len(*e) > 0 {
		return e
	}
	return nil
}
