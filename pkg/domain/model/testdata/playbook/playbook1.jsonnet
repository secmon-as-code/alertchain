{
  scenarios: [
    {
      name: 'test1',
      alert: import 'input.json',
      schema: 'scc',
      results: {
        ticket: [
          {
            name: 'test1',
          },
        ],
      },
    },
  ],
}
