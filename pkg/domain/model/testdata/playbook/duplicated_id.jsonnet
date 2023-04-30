// duplicated scenario id in playbook

{
  scenarios: [
    {
      id: 'test1',
      title: 'test1_title',
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
    {
      id: 'test1',  // duplicated
      title: 'test2_title',
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
