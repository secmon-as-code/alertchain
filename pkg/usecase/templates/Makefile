CMD_ALERTCHAIN=alertchain
CMD_OPA=opa

BASE_POLICY_FILES=\
	policy/alert/*.rego \
	policy/action/*.rego

TEST_POLICY_FILES=\
	policy/authz/*.rego \
	policy/play/*.rego

RESULT_FILE=policy/play/output/result.json

SCENARIO_FILES=\
	scenario/*.jsonnet \
	scenario/data/*.json

all: test

test: $(RESULT_FILE) $(BASE_POLICY_FILES) $(TEST_POLICY_FILES)
	$(CMD_OPA) test -v ./policy

$(RESULT_FILE): $(BASE_POLICY_FILES) $(SCENARIO_FILES)
	$(CMD_ALERTCHAIN) play -d ./policy -s ./scenario -o ./policy/play/output
