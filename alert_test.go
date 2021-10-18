package alertchain_test

/*

func setupAlertTest(t *testing.T, chain *alertchain.Chain) (usecase.Interface, infra.Clients, *types.Context) {
	clients := infra.Clients{
		DB: db.NewDBMock(t),
	}
	uc := usecase.New(clients, chain.Jobs.Convert(), chain.Actions.Convert())

	var wg sync.WaitGroup

	return uc, clients, types.NewContext().InjectWaitGroup(&wg)
}

type mock struct {
	Exec func(alert *alertchain.Alert) error
}

func (x *mock) Name() string { return "mock" }
func (x *mock) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	return x.Exec(alert)
}

func TestHandleAlert(t *testing.T) {
	var done bool
	var chain alertchain.Chain
	chain.NewJob().AddTask(&mock{
		Exec: func(alert *alertchain.Alert) error {
			alert.UpdateSeverity(types.SevAffected)
			alert.UpdateStatus(types.StatusClosed)
			done = true
			return nil
		},
	})
	uc, clients, ctx := setupAlertTest(t, &chain)

	input := ent.Alert{
		Title:    "five",
		Detector: "blue",
	}

	ctx.WaitGroup().Add(1)

	alert, err := uc.HandleAlert(ctx, &input, nil)
	require.NoError(t, err)
	require.NotNil(t, alert)

	ctx.WaitGroup().Wait()
	assert.True(t, done)

	got, err := clients.DB.GetAlert(ctx, alert.ID)
	require.NoError(t, err)
	assert.Equal(t, alert.Title, got.Title)
	assert.Equal(t, types.SevAffected, got.Severity)
	assert.Equal(t, types.StatusClosed, got.Status)
}

func TestRecvAlertDoNotUpdate(t *testing.T) {
	t.Run("do not update severity and status by overwriting vars", func(t *testing.T) {
		var done bool
		var chain alertchain.Chain
		chain.NewJob().AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error {
				alert.Severity = types.SevAffected
				alert.Status = types.StatusClosed
				done = true
				return nil
			},
		})
		uc, clients, ctx := setupAlertTest(t, &chain)

		input := alertchain.Alert{
			Alert: ent.Alert{
				Title:    "five",
				Detector: "blue",
			},
		}

		ctx.WaitGroup().Add(1)
		alert, err := uc.HandleAlert(ctx, &input.Alert, nil)
		ctx.WaitGroup().Wait()
		require.NoError(t, err)
		require.NotNil(t, alert)

		assert.True(t, done)

		got, err := clients.DB.GetAlert(ctx, alert.ID)
		require.NoError(t, err)
		assert.Equal(t, alert.Title, got.Title)
		assert.NotEqual(t, types.SevAffected, got.Severity)
		assert.NotEqual(t, types.StatusClosed, got.Status)
	})
}

func TestRecvAlertMassiveAnnotation(t *testing.T) {
	const multiplex = 32

	var chain alertchain.Chain
	job := chain.NewJob()
	job.Timeout = time.Second
	for i := 0; i < multiplex; i++ {
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error {
				require.Len(t, alert.Attributes, 1)
				alert.Attributes[0].Annotate(&alertchain.Annotation{
					Annotation: ent.Annotation{
						Source:    "x",
						Timestamp: rand.Int63(), // nosec
						Name:      "y",
						Value:     "z",
					},
				})
				return nil
			},
		})
	}
	uc, clients, ctx := setupAlertTest(t, &chain)

	inputAlert := ent.Alert{
		Title:    "five",
		Detector: "blue",
	}
	inputAttrs := []*ent.Attribute{
		{
			Key:   "color",
			Value: "red",
			Type:  types.AttrUserID,
		},
	}

	ctx.WaitGroup().Add(1)
	created, err := uc.HandleAlert(ctx, &inputAlert, inputAttrs)
	require.NoError(t, err)
	ctx.WaitGroup().Wait()

	alert, err := clients.DB.GetAlert(ctx, created.ID)
	require.NoError(t, err)
	require.Len(t, alert.Edges.Attributes[0].Edges.Annotations, multiplex)
	for _, ann := range alert.Edges.Attributes[0].Edges.Annotations {
		assert.Equal(t, "x", ann.Source)
		assert.Equal(t, "y", ann.Name)
		assert.Equal(t, "z", ann.Value)
		assert.Greater(t, ann.Timestamp, int64(0))
	}
}

func TestRecvAlertErrorHandling(t *testing.T) {
	t.Run("exit on error", func(t *testing.T) {
		var chain alertchain.Chain
		job := chain.NewJob()
		job.ExitOnErr = true
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error { return nil },
		})
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error { return errors.New("bomb!") },
		})

		done2ndJob := false
		chain.NewJob().AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error {
				done2ndJob = true
				return nil
			},
		})

		uc, _, ctx := setupAlertTest(t, &chain)

		input := ent.Alert{
			Title:    "five",
			Detector: "blue",
		}

		ctx.WaitGroup().Add(1)
		_, err := uc.HandleAlert(ctx, &input, nil)
		require.NoError(t, err)
		ctx.WaitGroup().Wait()
		assert.False(t, done2ndJob)
	})

	t.Run("not exit on error", func(t *testing.T) {
		var chain alertchain.Chain

		job := chain.NewJob()
		// Default: job.ExitOnErr = false
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error { return nil },
		})
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error { return errors.New("bomb!") },
		})

		done2ndJob := false
		chain.NewJob().AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error {
				done2ndJob = true
				return nil
			},
		})

		uc, _, ctx := setupAlertTest(t, &chain)

		input := ent.Alert{
			Title:    "five",
			Detector: "blue",
		}

		ctx.WaitGroup().Add(1)
		_, err := uc.HandleAlert(ctx, &input, nil)
		require.NoError(t, err)
		ctx.WaitGroup().Wait()
		assert.True(t, done2ndJob)
	})
}
*/
