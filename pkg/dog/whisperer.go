package dog

import (
	"context"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Whisperer manages cloud interactions associated with dogs.
// The Whisperer uses GCP Firestore as the data backend and GCP Cloud Tasks for scheduled jobs.
type Whisperer struct {
	QueueName  string
	TaskClient *cloudtasks.Client
	DogStore   DoggoStore
}

// DoggoStore is a data store for dogs.
type DoggoStore interface {
	Get(ctx context.Context, ID string) (*Dog, error)
	Put(ctx context.Context, dog *Dog) error
	Delete(ctx context.Context, ID string) error
}

// Register registers a dog by initializing its task and putting it in the data store.
// WARNING: on success, this method modifies the dog argument's NextTask fields.
func (w Whisperer) Register(ctx context.Context, dog *Dog) (*Dog, error) {
	// determine first scheduled time
	schedule, err := dog.Schedule()
	if err != nil {
		return nil, err
	}
	scheduleTime := schedule.Next(time.Now())

	// create task
	task, err := w.TaskClient.CreateTask(ctx, &tasks.CreateTaskRequest{
		Parent: w.QueueName,
		Task: &tasks.Task{
			MessageType: &tasks.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &tasks.AppEngineHttpRequest{
					HttpMethod:  tasks.HttpMethod_POST,
					RelativeUri: fmt.Sprintf("/dogs/%s/bark", dog.ID),
				},
			},
			ScheduleTime: timestamppb.New(scheduleTime),
		},
	})
	if err != nil {
		return dog, err
	}

	// update task fields of dog
	dog.NextTaskName = task.Name
	dog.NextTaskTime = task.ScheduleTime.AsTime()

	// insert into data store
	return dog, w.DogStore.Put(ctx, dog)
}

// Unregister deletes the dog and its associated task.
func (w Whisperer) Unregister(ctx context.Context, dogID string) error {
	// retrieve dog document
	dog, err := w.DogStore.Get(ctx, dogID)
	if err != nil {
		return err
	}

	// delete task
	err = w.TaskClient.DeleteTask(ctx, &tasks.DeleteTaskRequest{
		Name: dog.NextTaskName,
	})
	if err != nil {
		return err
	}

	// delete dog
	return w.DogStore.Delete(ctx, dogID)
}
