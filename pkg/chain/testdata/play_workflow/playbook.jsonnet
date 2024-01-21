{
  scenarios: [
    {
      id: 's1',
      title: 'Scenario 1',
      events: [
        {
          input: {
            class: 'threat',
          },
          schema: 'my_test',
          actions: {
            mock: [
              {
                index: 'first',
              },
              {
                index: 'second',
              },
            ],
          },
        },
      ],
    },
  ],
}
