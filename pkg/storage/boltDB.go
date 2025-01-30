package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/vvvdwbvvv/rover/pkg/model"

	"go.etcd.io/bbolt"
)

var (
	// 定義 BoltDB 存儲的 Bucket 名稱
	containerBucket = []byte("containers")

	// 常見錯誤
	ErrBucketNotFound    = errors.New("storage bucket not found")
	ErrContainerNotFound = errors.New("container not found")
)

// BoltDB 存儲管理
type BoltDB struct {
	db *bbolt.DB
}

// NewBoltDB 初始化 BoltDB
func NewBoltDB(dbPath string) (*BoltDB, error) {
	options := &bbolt.Options{
		Timeout: 1 * time.Second, // 超時時間，防止長時間鎖定
		NoSync:  false,           // 確保數據持久化
	}

	db, err := bbolt.Open(dbPath, 0600, options)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 確保 "containers" bucket 存在
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(containerBucket)
		return err
	}); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return &BoltDB{db: db}, nil
}

// SaveContainer 存儲容器狀態
func (b *BoltDB) SaveContainer(container model.ContainerState) error {
	if container.Name == "" {
		return errors.New("container name cannot be empty")
	}

	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(containerBucket)
		if bucket == nil {
			return ErrBucketNotFound
		}

		data, err := json.Marshal(container)
		if err != nil {
			return fmt.Errorf("failed to marshal container: %w", err)
		}

		return bucket.Put([]byte(container.Name), data)
	})
}

// GetContainer 取得單個容器
func (b *BoltDB) GetContainer(name string) (*model.ContainerState, error) {
	var container model.ContainerState

	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(containerBucket)
		if bucket == nil {
			return ErrBucketNotFound
		}

		data := bucket.Get([]byte(name))
		if data == nil {
			return ErrContainerNotFound
		}

		return json.Unmarshal(data, &container)
	})

	if err != nil {
		return nil, err
	}
	return &container, nil
}

// GetContainers 取得所有容器
func (b *BoltDB) GetContainers() ([]model.ContainerState, error) {
	var containers []model.ContainerState

	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(containerBucket)
		if bucket == nil {
			return ErrBucketNotFound
		}

		return bucket.ForEach(func(k, v []byte) error {
			var container model.ContainerState
			if err := json.Unmarshal(v, &container); err != nil {
				return fmt.Errorf("failed to unmarshal container %s: %w", k, err)
			}
			containers = append(containers, container)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return containers, nil
}

// DeleteContainer 刪除容器
func (b *BoltDB) DeleteContainer(name string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(containerBucket)
		if bucket == nil {
			return ErrBucketNotFound
		}

		return bucket.Delete([]byte(name)) // 直接刪除，不需要先 Get()
	})
}

// Close 關閉 BoltDB
func (b *BoltDB) Close() error {
	if b.db == nil {
		return nil
	}
	return b.db.Close()
}
