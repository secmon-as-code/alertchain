package usecase

import (
	"bytes"
	"context"
	"os"
	"sync"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func (x *usecase) GetAlerts(ctx *types.Context) ([]*ent.Alert, error) {
	return x.clients.DB.GetAlerts(ctx)
}

func (x *usecase) GetAlert(ctx *types.Context, id types.AlertID) (*ent.Alert, error) {
	return x.clients.DB.GetAlert(ctx, id)
}

func (x *usecase) HandleAlert(ctx *types.Context, alert *ent.Alert, attrs []*ent.Attribute) (*ent.Alert, error) {
	created, err := saveAlert(ctx, x.clients.DB, alert, attrs)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := executeJobs(ctx, &x.clients, x.jobs, created.ID); err != nil {
			utils.HandleError(err)
		}
		if wg := ctx.WaitGroup(); wg != nil {
			logger.Trace("Done waitgroup")
			wg.Done()
		}
	}()

	return created, nil
}

func saveAlert(ctx *types.Context, dbClient db.Interface, recv *ent.Alert, attrs []*ent.Attribute) (*ent.Alert, error) {
	if err := ValidateAlert(recv); err != nil {
		return nil, goerr.Wrap(err)
	}

	created, err := dbClient.NewAlert(ctx)
	if err != nil {
		return nil, err
	}

	if err := dbClient.UpdateAlert(ctx, created.ID, recv); err != nil {
		return nil, err
	}
	if err := dbClient.AddAttributes(ctx, created.ID, attrs); err != nil {
		return nil, err
	}

	newAlert, err := dbClient.GetAlert(ctx, created.ID)
	if err != nil {
		return nil, err
	}

	return newAlert, nil
}

func executeJobs(ctx *types.Context, clients *infra.Clients, jobs []*Job, alertID types.AlertID) error {
	logger.With("alertID", alertID).Trace("Starting executeJobs")

	for idx, job := range jobs {
		logger.With("job", job).Trace("Starting Job")

		if len(job.Tasks) == 0 {
			continue
		}
		if job.Timeout > 0 {
			newCtx, cancel := context.WithTimeout(ctx, job.Timeout)
			defer cancel()
			ctx = types.WrapContext(newCtx)
		}

		alert, err := clients.DB.GetAlert(ctx, alertID)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		errCh := make(chan error, len(job.Tasks))

		for _, task := range job.Tasks {
			wg.Add(1)
			go executeTask(ctx, &executeTaskInput{
				stage:  int64(idx),
				task:   task,
				wg:     &wg,
				alert:  alert,
				client: clients.DB,
				errCh:  errCh,
			})
		}
		wg.Wait()
		logger.With("job", job).Trace("Completed Job")

		close(errCh)
		for err := range errCh {
			utils.HandleError(err)
			if err != nil && job.ExitOnErr {
				return err
			}
		}
	}

	logger.With("alertID", alertID).Trace("Exiting executeJobs")

	return nil
}

type executeTaskInput struct {
	stage  int64
	task   *Task
	wg     *sync.WaitGroup
	alert  *ent.Alert
	client db.Interface
	errCh  chan error
}

type logWriter struct {
	bytes.Buffer
}

func (x *logWriter) Write(p []byte) (n int, err error) {
	if n, err := x.Buffer.Write(p); err != nil {
		return n, err
	}
	return os.Stdout.Write(p)
}

func executeTask(ctx *types.Context, input *executeTaskInput) {
	defer input.wg.Done()
	var taskLog *ent.TaskLog
	var execLog ent.ExecLog
	var logW logWriter
	ctx = ctx.InjectWriter(&logW)

	handlerError := func(err error) {
		wrapped := goerr.Wrap(err).With("task.Name", input.task.Name)
		utils.HandleError(wrapped)
		input.errCh <- wrapped

		if taskLog != nil {
			utils.CopyErrorToExecLog(err, &execLog)
		}
	}

	defer func() {
		if taskLog != nil {
			execLog.Log = logW.String()
			execLog.Timestamp = time.Now().UTC().UnixNano()

			if err := input.client.AppendTaskLog(ctx, taskLog.ID, &execLog); err != nil {
				input.errCh <- err
			}
		}
	}()

	taskLog, err := input.client.NewTaskLog(ctx, input.alert.ID, input.task.Name, input.stage)
	if err != nil {
		handlerError(err)
		return
	}

	if err := input.client.AppendTaskLog(ctx, taskLog.ID, &ent.ExecLog{
		Status: types.ExecStart,
	}); err != nil {
		handlerError(err)
		return
	}

	req, err := input.task.Execute(ctx, input.alert)
	if err != nil {
		handlerError(err)
		return
	}

	if err := commitAlert(ctx, input.client, input.alert.ID, req); err != nil {
		handlerError(err)
		return
	}
}

func commitAlert(ctx *types.Context, client db.Interface, id types.AlertID, req *ChangeRequest) error {
	ts := time.Now().UTC().Unix()
	if req.newStatus != nil {
		if err := client.UpdateAlertStatus(ctx, id, *req.newStatus, ts); err != nil {
			return err
		}
	}
	if req.newSeverity != nil {
		if err := client.UpdateAlertSeverity(ctx, id, *req.newSeverity, ts); err != nil {
			return err
		}
	}

	if len(req.newAttrs) > 0 {
		if err := client.AddAttributes(ctx, id, req.newAttrs); err != nil {
			return err
		}
	}

	for _, newAnn := range req.newAnnotations {
		if err := client.AddAnnotation(ctx, newAnn.attr, []*ent.Annotation{newAnn.ann}); err != nil {
			return err
		}
	}

	for _, ref := range req.newReferences {
		if err := client.AddReference(ctx, id, ref); err != nil {
			return err
		}
	}

	return nil
}
