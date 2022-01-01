package model_test

/*
type PushAlert struct {
	err error
}

func (x *PushAlert) Name() string { return "blue" }
func (x *PushAlert) Run(handler model.Handler) error {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		err := handler(ctx, &model.Alert{
			Title:    fmt.Sprintf("%d", i),
			Detector: "blue",
		})
		if err != nil {
			x.err = err
			break
		}
	}
	return nil
}

type fallback struct {
	req *http.Request
}

func (x *fallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x.req = r
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("good"))
}

func bind(t *testing.T, body io.Reader, dst interface{}) {
	raw, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, dst))
}

func TestChainWithAPI(t *testing.T) {
	addr := "127.0.0.1:45678"
	mock := db.NewDBMock(t)
	src := &PushAlert{}
	chain, err := alertchain.New(
		alertchain.WithDB(mock),
		alertchain.WithSources(src),
		alertchain.WithAPI(addr, "https://alertchain.example.com/", nil),
	)
	require.NoError(t, err)
	go func() {
		chain.Start()
	}()
	time.Sleep(time.Millisecond * 500)
	resp, err := http.Get("http://" + addr + "/api/v1/alert")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var alerts []*alertchain.Alert
	bind(t, resp.Body, &alerts)
	require.Len(t, alerts, 10)
	assert.Equal(t, "0", alerts[0].Title)
	assert.Equal(t, "9", alerts[9].Title)

	require.Len(t, alerts[0].References, 1)
	assert.Equal(t, fmt.Sprintf("https://alertchain.example.com/api/v1/alert/%s", alerts[0].ID),
		alerts[0].References[0].URL)
}

func TestAlertAPI(t *testing.T) {
	mock := db.NewDBMock(t)
	src := &PushAlert{}
	chain, err := alertchain.New(
		alertchain.WithDB(mock),
		alertchain.WithSources(src),
	)
	require.NoError(t, err)
	require.NoError(t, chain.Start())
	require.NoError(t, src.err)

	logger := zlog.New()

	t.Run("without fallback", func(t *testing.T) {
		engine := alertchain.NewAPIEngine(mock, nil, logger)

		t.Run("get alerts", func(t *testing.T) {
			var id types.AlertID
			{
				w := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "http://example.com/api/v1/alert", nil)
				require.NoError(t, err)
				engine.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)

				var alerts []*alertchain.Alert
				bind(t, w.Result().Body, &alerts)
				require.Len(t, alerts, 10)
				assert.Equal(t, "0", alerts[0].Title)
				assert.Equal(t, "9", alerts[9].Title)
				id = alerts[0].ID
			}

			{
				w := httptest.NewRecorder()
				url := "http://example.com/api/v1/alert/" + string(id)
				t.Log(url)
				req, err := http.NewRequest("GET", url, nil)
				t.Log(id)
				require.NoError(t, err)
				engine.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)

				var alert alertchain.Alert
				bind(t, w.Result().Body, &alert)
				assert.Equal(t, "0", alert.Title)
			}
		})

		t.Run("get alert", func(t *testing.T) {
			{
				w := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "http://example.com/api/v1/alert", nil)
				require.NoError(t, err)
				engine.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)

				var alerts []*alertchain.Alert
				bind(t, w.Result().Body, &alerts)
				require.Len(t, alerts, 10)
				assert.Equal(t, "0", alerts[0].Title)
				assert.Equal(t, "9", alerts[9].Title)
			}
		})

		t.Run("not found", func(t *testing.T) {
			{
				w := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "http://example.com/xxx", nil)
				require.NoError(t, err)
				engine.ServeHTTP(w, req)
				assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
			}
		})
	})

	t.Run("with fallback", func(t *testing.T) {
		f := &fallback{}
		engine := alertchain.NewAPIEngine(mock, f, logger)

		t.Run("get alerts", func(t *testing.T) {
			{
				w := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "http://example.com/api/v1/alert", nil)
				require.NoError(t, err)
				engine.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)

				var alerts []*alertchain.Alert
				bind(t, w.Result().Body, &alerts)
				require.Len(t, alerts, 10)
			}
		})

		t.Run("to fallback", func(t *testing.T) {
			{
				w := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "http://example.com/xxx", nil)
				require.NoError(t, err)
				engine.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)

				raw, err := ioutil.ReadAll(w.Result().Body)
				require.NoError(t, err)
				assert.Equal(t, "good", string(raw))

				require.NotNil(t, f.req)
				assert.Equal(t, "/xxx", f.req.URL.Path)
			}
		})
	})
}
*/
