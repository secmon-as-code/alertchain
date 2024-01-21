{
  id: 'test1',
  title: 'test1_title',
  events: [
    {
      input: import 'input.json',
      schema: 'scc',
      actions: {
        'chatgpt.query': [
          {
            name: 'test1',
          },
        ],
      },
    },
  ],
}
