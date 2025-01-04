local event = import 'data/event.json';

{
  id: 'my_first_scenario',
  title: 'Create an issue on GitHub and add a comment to it',
  events: [
    # The first alert should create an issue
    {
      input: event,
      schema: 'your_schema',
      actions: {
        "github.create_issue": [{number: 666}],
      },
    },

    # The second alert should add a comment to the issue
    {
      input: event,
      schema: 'your_schema',
    },
  ],

  env: import 'env.libsonnet',
}
