{
  policy: {},
  actions: [
    {
      id: 'test-scc',
      uses: 'scc',
      config: {
        data: std.extVar('COLOR'),
        num: 123,
      },
    },
  ],
}
