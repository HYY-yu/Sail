package svc_publish

import (
	"context"
	"fmt"
	"github.com/HYY-yu/sail/internal/service/sail/api/repo"
	"github.com/HYY-yu/sail/internal/service/sail/model"
	"github.com/HYY-yu/sail/internal/service/sail/storage"
	"github.com/HYY-yu/seckill.pkg/db"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestPublishSvc_initPublish(t *testing.T) {
	etcdStorage, err := storage.New(&storage.ETCDConfig{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	assert.NoError(t, err)
	ctrl := gomock.NewController(t)
	cs := NewMockConfigSystem(ctrl)
	publishMgr := repo.NewMockPublishMgrInter(ctrl)
	publishMgr.EXPECT().CreatePublish(gomock.Any()).Return(nil).AnyTimes()

	publishRepo := repo.NewMockPublishRepo(ctrl)
	publishRepo.EXPECT().Mgr(gomock.Any(), gomock.Any()).Return(publishMgr).AnyTimes()

	// 替换 ConfigSystem
	cs.EXPECT().GetConfigProjectAndNamespace(context.Background(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, projectID int, namespaceID int) (*model.Project, *model.Namespace, error) {
			return &model.Project{ID: projectID, ProjectGroupID: 1, Key: fmt.Sprintf("TEST%d", projectID), Name: "TEST"},
				&model.Namespace{ID: namespaceID, ProjectGroupID: 1, Name: fmt.Sprintf("TEST%d", namespaceID), SecretKey: fmt.Sprintf("TEST%d", namespaceID)}, nil
		}).AnyTimes()

	dbRepo := db.MockRepo{}

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
				DB:           dbRepo,
				Store:        etcdStorage,
				PublishRepo:  publishRepo,
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
				for i := range tt.args {
					tokens = append(tokens, <-tokenChan)
					_ = i
				}
			} else {
				for i := range tt.args {
					token, err := publishSvc.initPublish(context.Background(), tt.args[i].projectID, tt.args[i].namespaceID)
					assert.NoError(t, err)
					tokens = append(tokens, token)
				}
			}
			t.Log(tokens)
			assert.Equal(t, tt.wantSame, isSame(tokens))
			// Clear
			for i := range tt.args {
				err := publishSvc.deletePublish(context.Background(), tt.args[i].projectID, tt.args[i].namespaceID)
				assert.NoError(t, err)
			}
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
