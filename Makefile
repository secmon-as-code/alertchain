all: gql mock

# ---------------------------
gql: pkg/controller/graphql/generated.go

pkg/controller/graphql/generated.go: graphql/schema.graphqls
	go run github.com/99designs/gqlgen@v0.17.61 generate

# ---------------------------
MOCK_OUT=pkg/mock/infra.go
MOCK_SRC=./pkg/domain/interfaces
MOCK_INTERFACES=GenAI Database

mock: $(MOCK_OUT)

$(MOCK_OUT): $(MOCK_SRC)/*
	go run github.com/matryer/moq@v0.5.1 -pkg mock -out $(MOCK_OUT) $(MOCK_SRC) $(MOCK_INTERFACES)

clean:
	rm -f $(MOCK_OUT)
