package bucket

import (
	"context"
	"fmt"
	"log"

	"io/ioutil"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"github.com/lokmannicholas/firego"
)

var ctx = context.Background()

type FireBucket interface {
	Create() error
	Delete(fileName string) error
	CreateFile(mimeType string, fireMeta map[string]string, fileName string, content []byte) error
	ReadFile(fileName string) error
	ListBucket(queryPrefix string) error
}

type FireBucketImpl struct {
	FirebaseStorage *storage.BucketHandle
}

func (b *FireBucketImpl) Create() error {

	attr := &storage.BucketAttrs{}
	err := b.FirebaseStorage.Create(ctx, firebase.GetConfig().ProjectID, attr)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (b *FireBucketImpl) Delete(fileName string) error {

	err := b.FirebaseStorage.Object(fileName).Delete(ctx)
	if err != nil {
		_ = fmt.Errorf("%+v", err)
		return err
	}
	return nil
}

func (b *FireBucketImpl) CreateFile(mimeType string, fireMeta map[string]string, fileName string, content []byte) error {
	wc := b.FirebaseStorage.Object(fileName).NewWriter(ctx)
	wc.ContentType = mimeType // "text/plain"
	wc.Metadata = fireMeta

	if _, err := wc.Write(content); err != nil {
		_ = fmt.Errorf("%+v", err)
		return err
	}
	if err := wc.Close(); err != nil {
		_ = fmt.Errorf("%+v", err)
		return err
	}
	return nil
}
func (b *FireBucketImpl) ReadFile(fileName string) error {

	rc, err := b.FirebaseStorage.Object(fileName).NewReader(ctx)
	if err != nil {
		_ = fmt.Errorf("%+v", err)
		return err
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		_ = fmt.Errorf("%+v", err)
		return err
	}

	if len(slurp) > 1024 {
		fmt.Println(slurp)
	} else {
	}
	return nil
}
func (b *FireBucketImpl) ListBucket(queryPrefix string) error {

	query := &storage.Query{Prefix: queryPrefix}
	it := b.FirebaseStorage.Objects(ctx, query)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(obj)
	}
	return nil
}

func GetFireBucket() FireBucket {
	//storage bucket
	client, err := firebase.GetFireApp().Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	firebucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}
	return &FireBucketImpl{
		FirebaseStorage: firebucket,
	}
}
