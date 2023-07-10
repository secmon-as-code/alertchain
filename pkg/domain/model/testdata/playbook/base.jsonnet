{
  scenarios: [
    {
      id: 'test1',
      title: 'test1_title',
      events: [
        {
          input: import 'input.json',
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
    },
  ],
}
