# Scenario Test

This directory contains the scenario test for alertchain. Scenario test can be done by `play` command. The scenario test is a test that runs with `alert` and `action` policies and a scenario file as JSONNET format.

## Directory Structure

- `event`: Event data for scenario test
  - `guardduty.json`: GuardDuty event data
- `policy`: Policy data for scenario test
  - `alert.json`: Alert policy data
  - `action.json`: Action policy data
- `scenario`: Scenario data for scenario test
  - `scenario.jsonnet`: Scenario data in JSONNET format
- `results`: The respond data of each action.
  - `chatgpt.json`: The respond data of ChatGPT action
- `output/scenario1/data.json`: The output data of scenario test (that will be created after running the scenario test)
- `test.rego`: Rego policy for testing the output data of scenario test

## How to Run

1. Run the scenario test by `play` command
    ```bash
    # at root directory of the repository
    $ alertchain play \
        -s ./examples/test/scenarios/scenario.jsonnet \
        -d ./examples/test/policy \
        -o ./examples/test/output
    ```
    - `-s`: The path of the scenario file
    - `-d`: The path of the policy directory
    - `-o`: The path of the output directory
2. Check the output data
    ```bash
    $ cat examples/test/output/scenario1/data.json | jq
3. Test the output data
    ```bash
    $ opa test -v test.rego examples/test/output/scenario1/data.json
    ```