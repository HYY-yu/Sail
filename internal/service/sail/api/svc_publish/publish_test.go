package svc_publish

import (
	"context"
	"fmt"
	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"github.com/HYY-yu/sail/internal/service/sail/storage"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"sync"
	"testing"
)

func TestPublishSvc_initPublish(t *testing.T) {
	etcdStorage, err := storage.New(&storage.ETCDConfig{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	assert.NoError(t, err)
	mockRepo := repo.NewMockPublishRepo()
	ctrl := gomock.NewController(t)
	cs := NewMockConfigSystem(ctrl)

	// 替换 ConfigSystem
	cs.EXPECT().GetConfigProjectAndNamespace(context.Background(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, projectID int, namespaceID int) (*model.Project, *model.Namespace, error) {
			return &model.Project{ID: projectID, ProjectGroupID: 1, Key: fmt.Sprintf("TEST%d", projectID), Name: "TEST"},
				&model.Namespace{ID: namespaceID, ProjectGroupID: 1, Name: fmt.Sprintf("TEST%d", namespaceID), SecretKey: fmt.Sprintf("TEST%d", namespaceID)}, nil
		})

	// 替换 pm
	pm := mockRepo.Mgr(context.Background(), &gorm.DB{})
	gomonkey.ApplyMethodFunc(pm, "CreatePublish", func(bean *model.Publish) (err error) {
		t.Log("Save Config Success!! ")
		return nil
	})

	type args struct {
		projectID   int
		namespaceID int
	}

	tests := []struct {
		name        string
		concurrency bool
		args        []args
		wantSame    bool
	}{
		{
			name: "TEST：顺序调用",
			args: []args{
				{
					projectID:   1,
					namespaceID: 1,
				},
				{
					projectID:   1,
					namespaceID: 1,
				},
			},
			wantSame: true,
		},
		{
			name: "TEST：顺序调用2",
			args: []args{
				{
					projectID:   1,
					namespaceID: 1,
				},
				{
					projectID:   2,
					namespaceID: 1,
				},
			},
			wantSame: false,
		},
		{
			name:        "TEST：并发调用",
			concurrency: true,
			args: []args{
				{
					projectID:   1,
					namespaceID: 1,
				},
				{
					projectID:   1,
					namespaceID: 1,
				},
			},
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publishSvc := &PublishSvc{
				Store:        etcdStorage,
				PublishRepo:  mockRepo,
				configSystem: cs,
			}
			tokens := make([]string, 0)
			if tt.concurrency {
				tokenChan := make(chan string, len(tt.args))
				defer close(tokenChan)
				wg := sync.WaitGroup{}
				for i := range tt.args {
					wg.Add(1)

					go func(i int) {
						token, err := publishSvc.initPublish(context.Background(), tt.args[i].projectID, tt.args[i].namespaceID)
						assert.NoError(t, err)
						tokenChan <- token
						wg.Done()
					}(i)
				}

				wg.Wait()
				for token := range tokenChan {
					tokens = append(tokens, token)
				}
			} else {
				for i := range tt.args {
					token, err := publishSvc.initPublish(context.Background(), tt.args[i].projectID, tt.args[i].namespaceID)
					assert.NoError(t, err)
					tokens = append(tokens, token)
				}
			}
			assert.Equal(t, tt.wantSame, isSame(tokens))
		})
	}
}

func isSame(arr []string) bool {
	if len(arr) == 1 || len(arr) == 0 {
		return true
	}

	for i := 0; i < len(arr)-1; i++ {
		if arr[i] != arr[i+1] {
			return false
		}
	}
	return true
}
