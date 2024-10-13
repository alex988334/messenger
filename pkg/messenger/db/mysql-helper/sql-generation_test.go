package mysql_helper

import (
	"testing"
)

func TestGenerateSelect(t *testing.T) {

	t.Log("Generate Select to sql")

	required := []string{
		SELECT + "cu.id_chat AS \"cu.id_chat\", cu.id_user AS \"cu.id_user\", " +
			"cu.client_hash AS \"cu.client_hash\", c.id AS \"c.id\", c.autor AS \"c.autor\", " +
			"c.status AS \"c.status\"",

		SELECT + "id_chat, id_user, client_hash, id, autor, status",
	}
	params := [][][]string{
		{
			{"id_chat", "id_user", "", "client_hash", "id", "autor", "status", "", ""},
			{"cu", "cu", "", "cu", "c", "c", "c", "", ""},
		},
		{
			{"id_chat", "id_user", "", "client_hash", "id", "autor", "status", "", ""},
			{},
		},
	}

	model := &SQLModel{}

	for i := 0; i < len(params); i++ {

		model.selectFields = params[i][0]
		model.selectAliases = params[i][1]

		sel := generateSelect(model)
		if sel != required[i] {
			t.Fatal("Failed generate Select to sql! Required:", required[i], ", rezult:", sel)
		}
	}
}

func TestGenerateFrom(t *testing.T) {

	t.Log("Generate From to sql")

	required := []string{
		FROM + "chat_user cu, chat c",
		FROM + "chat_user",
	}
	params := [][][]string{
		{
			{"cu", "", "c", ""},
			{"chat_user", "", "chat", ""},
		},
		{
			{},
			{"chat_user"},
		},
	}
	model := &SQLModel{}

	for i := 0; i < len(params); i++ {

		model.sqlTablesAliases = params[i][0]
		model.sqlTables = params[i][1]

		fro := generateFrom(model)
		if fro != required[i] {
			t.Fatal("Failed generate From to sql! Required:", required[i], ", rezult:", fro)
		}
	}
}

func TestGenerateWhereDependencies(t *testing.T) {

	t.Log("Generate Where Dependencies to sql")

	model := &SQLModel{}
	required := WHERE + "cu.id_user" + EQUAL + "c.autor" + AND + "c.id" + EQUAL + "m.id_chat"

	model.whereDependenciesAliases = []string{
		"cu", "c", "c", "m",
	}
	model.whereDependenciesFields = []string{
		"id_user", "autor", "id", "id_chat",
	}
	model.whereDependenciesOperator = []string{EQUAL, ""}
	model.whereDependenciesLogicOperators = []string{}

	dep := generateWhereDependencies(model)
	if dep != required {
		t.Fatal("Failed generate Where Dependencies to sql! Required:", required, ", rezult:", dep)
	}
}

func TestGenerateWhere(t *testing.T) {

	t.Log("Generate Where to sql")

	model := &SQLModel{}

	required := []string{
		WHERE + "cu.id_user" + EQUAL + PARAMETR,

		WHERE + "cu.id_user" + EQUAL + PARAMETR + AND + "(c.autor" + IS + NULL + OR +
			"c.id" + NOT_EQUAL + PARAMETR + ")" + AND + "id_chat" + MORE + PARAMETR,

		WHERE + "(cu.id_user" + EQUAL + PARAMETR + OR + "c.autor" + IS + NOT + NULL + OR +
			"c.id" + NOT_EQUAL + PARAMETR + ")" + AND + "id_chat" + MORE + PARAMETR,
	}
	params := [][][]string{
		{
			{"cu"},
			{"id_user"},
			{EQUAL},
			{},
		},
		{
			{"cu", "c", "c", "", "", ""},
			{"id_user", "autor", "id", "id_chat"},
			{EQUAL, NULL, NOT_EQUAL, MORE, "", ""},
			{AND, OR, AND, AND},
		},
		{
			{"cu", "c", "c", "", "", ""},
			{"id_user", "autor", "id", "id_chat"},
			{EQUAL, NOT + NULL, NOT_EQUAL, MORE, "", ""},
			{OR, OR, AND, AND},
		},
	}

	for i := 0; i < len(params); i++ {

		model.whereAliases = params[i][0]
		model.whereFields = params[i][1]
		model.whereFieldsOperator = params[i][2]
		model.whereLogicOperators = params[i][3]

		wher := generateWhereConditions(model)
		if wher != required[i] {
			t.Fatal("Failed generate Where to sql! \nRequired:", required[i], "; \nrezult:  ", wher)
		}
	}
}

func TestGenerateJoin(t *testing.T) {

	t.Log("Generate Join to sql")

	model := &SQLModel{}

	required := []string{
		LEFT + JOIN + "chat_user cu" + ON + "c.author" + EQUAL + "cu.id_user" +
			JOIN + "user u" + ON + "cu.id_user" + EQUAL + "u.id",

		JOIN + "chat_user cu" + ON + "c.author" + EQUAL + "cu.id_user" +
			JOIN + "user u" + ON + "cu.id_user" + EQUAL + "u.id",
	}
	params := [][][]string{
		{
			{"cu", "u"},
			{LEFT, ""},
			{"c", "cu", "cu", "u"},
			{"author", "id_user", "id_user", "id"},
			{"chat_user", "user"},
		},
		{
			{"cu", "u"},
			{"", ""},
			{"c", "cu", "cu", "u"},
			{"author", "id_user", "id_user", "id"},
			{"chat_user", "user"},
		},
	}

	for i := 0; i < len(params); i++ {

		model.joinAliases = params[i][0]
		model.joinDirection = params[i][1]
		model.joinLinkAlias = params[i][2]
		model.joinLinkFields = params[i][3]
		model.joinTables = params[i][4]

		rezult := generateJoin(model)
		if rezult != required[i] {
			t.Fatal("Failed generate Join to sql! \nRequired:", required[i], "; \nrezult:  ", rezult)
		}
	}
}

func TestGenerateGroupBy(t *testing.T) {

	t.Log("Generate Group By to sql")

	model := &SQLModel{}

	required := []string{
		GROUP_BY + "cu.id_user" + DELIMETR + "cu.id_chat",
	}
	params := [][][]string{
		{
			{"cu", "cu"},
			{"id_user", "id_chat"},
		},
	}

	for i := 0; i < len(params); i++ {

		model.groupAliases = params[i][0]
		model.groupFields = params[i][1]

		rezult := generateGroupBy(model)
		if rezult != required[i] {
			t.Fatal("Failed generate Group By to sql! \nRequired:", required[i], "; \nrezult:  ", rezult)
		}
	}
}

func TestGenerateOrderBy(t *testing.T) {

	t.Log("Generate Order By to sql")

	model := &SQLModel{}

	required := []string{
		ORDER_BY + "cu.id_user" + ASC + DELIMETR + "cu.id_chat" + DESC,
	}
	params := [][][]string{
		{
			{"cu", "cu"},
			{"id_user", "id_chat"},
			{ASC, DESC},
		},
	}

	for i := 0; i < len(params); i++ {

		model.orderAliases = params[i][0]
		model.orderFileds = params[i][1]
		model.orderDirection = params[i][2]

		rezult := generateOrderBy(model)
		if rezult != required[i] {
			t.Fatal("Failed generate Order By to sql! \nRequired:", required[i], "; \nrezult:  ", rezult)
		}
	}
}

func TestGenerateLimit(t *testing.T) {

	t.Log("Generate Limit to sql")

	model := &SQLModel{}

	required := []string{
		LIMIT + "25",
		LIMIT + "1",
	}
	params := []int{25, 0}

	for i := 0; i < len(params); i++ {
		model.limit = params[i]
		rezult := generateLimit(model)
		if rezult != required[i] {
			t.Fatal("Failed generate Limit to sql! \nRequired:", required[i], "; \nrezult:  ", rezult)
		}
	}
}

func TestPrepareSelectQuery(t *testing.T) {

	t.Log("Generate Prepare Select Query to sql")

	var model *SQLModel

	required := []string{
		"SELECT id, author, alias, create_at, status FROM chat WHERE id=?",

		"SELECT c.id AS \"c.id\", c.author AS \"c.author\", c.alias AS \"c.alias\", " +
			"c.create_at AS \"c.create_at\", c.status AS \"c.status\" FROM chat c WHERE c.id=?",

		"SELECT c.id AS \"c.id\", c.author AS \"c.author\", c.alias AS \"c.alias\", " +
			"u.username AS \"u.username\" FROM chat c, user u, chat_user cu " +
			"WHERE c.id=u.author AND c.id=cu.id_chat AND c.id=?",

		"SELECT c.id AS \"c.id\", c.author AS \"c.author\", c.alias AS \"c.alias\", " +
			"u.username AS \"u.username\" FROM chat c LEFT JOIN user u ON c.id=u.author" +
			" JOIN chat_user cu ON c.id=cu.id_chat WHERE c.id=?",

		"SELECT c.id AS \"c.id\", c.author AS \"c.author\", c.alias AS \"c.alias\", " +
			"u.username AS \"u.username\"" +
			" FROM chat c LEFT JOIN user u ON c.author=u.id" +
			" JOIN chat_user cu ON c.id=cu.id_chat WHERE c.id=? " +
			"GROUP BY c.author, cu.id_chat ORDER BY c.author ASC, cu.id_chat DESC LIMIT 24",
	}
	params := []*SQLModel{
		&SQLModel{
			selectFields:        []string{"id", "author", "alias", "create_at", "status"},
			selectAliases:       []string{},
			sqlTables:           []string{"chat"},
			sqlTablesAliases:    []string{},
			whereFields:         []string{"id"},
			whereAliases:        []string{},
			whereFieldsOperator: []string{EQUAL},
			whereLogicOperators: []string{},
		},
		&SQLModel{
			selectFields:        []string{"id", "author", "alias", "create_at", "status"},
			selectAliases:       []string{"c", "c", "c", "c", "c"},
			sqlTables:           []string{"chat"},
			sqlTablesAliases:    []string{"c"},
			whereFields:         []string{"id"},
			whereAliases:        []string{"c"},
			whereFieldsOperator: []string{EQUAL},
			whereLogicOperators: []string{},
		},
		&SQLModel{
			selectFields:                    []string{"id", "author", "alias", "username"},
			selectAliases:                   []string{"c", "c", "c", "u"},
			sqlTables:                       []string{"chat", "user", "chat_user"},
			sqlTablesAliases:                []string{"c", "u", "cu"},
			whereFields:                     []string{"id"},
			whereAliases:                    []string{"c"},
			whereFieldsOperator:             []string{EQUAL},
			whereLogicOperators:             []string{},
			whereDependenciesFields:         []string{"id", "author", "id", "id_chat"},
			whereDependenciesAliases:        []string{"c", "u", "c", "cu"},
			whereDependenciesOperator:       []string{EQUAL, EQUAL},
			whereDependenciesLogicOperators: []string{AND},
		},
		&SQLModel{
			selectFields:        []string{"id", "author", "alias", "username"},
			selectAliases:       []string{"c", "c", "c", "u"},
			sqlTables:           []string{"chat"},
			sqlTablesAliases:    []string{"c", "u", "cu"},
			whereFields:         []string{"id"},
			whereAliases:        []string{"c"},
			whereFieldsOperator: []string{EQUAL},
			whereLogicOperators: []string{},
			joinTables:          []string{"user", "chat_user"},
			joinAliases:         []string{"u", "cu"},
			joinDirection:       []string{LEFT, ""},
			joinLinkFields:      []string{"id", "author", "id", "id_chat"},
			joinLinkAlias:       []string{"c", "u", "c", "cu"},
		},
		&SQLModel{
			selectFields:        []string{"id", "author", "alias", "username"},
			selectAliases:       []string{"c", "c", "c", "u"},
			sqlTables:           []string{"chat"},
			sqlTablesAliases:    []string{"c"},
			whereFields:         []string{"id"},
			whereAliases:        []string{"c"},
			whereFieldsOperator: []string{EQUAL},
			whereLogicOperators: []string{},
			joinTables:          []string{"user", "chat_user"},
			joinAliases:         []string{"u", "cu"},
			joinDirection:       []string{LEFT, ""},
			joinLinkFields:      []string{"author", "id", "id", "id_chat"},
			joinLinkAlias:       []string{"c", "u", "c", "cu"},
			groupFields:         []string{"author", "id_chat"},
			groupAliases:        []string{"c", "cu"},
			orderFileds:         []string{"author", "id_chat"},
			orderAliases:        []string{"c", "cu"},
			orderDirection:      []string{ASC, DESC},
			limit:               24,
		},
	}

	for i := 0; i < len(params); i++ {

		model = params[i]

		model.PrepareQuery()

		if model.query != required[i] {
			t.Fatal("Failed generate Prepare Query to sql! \nRequired:", required[i], "; \nrezult:  ", model.query)
		}
	}
}

func TestPrepareUpdateQuery(t *testing.T) {

	t.Log("Generate Prepare Update Query to sql")

	var model *SQLModel

	required := []string{
		"UPDATE chat SET alias=?, create_at=? WHERE id=?",

		"UPDATE chat c SET alias=?, c.create_at=? WHERE c.id=? AND author=?",

		"UPDATE chat c, chat_user cu, user u SET c.alias=?, c.create_at=?, cu.client_hash=? " +
			"WHERE c.id=cu.id_chat AND u.id=c.author AND c.id=? AND c.author=?",
	}
	params := []*SQLModel{
		&SQLModel{
			updSetFields:        []string{"alias", "create_at"},
			updSetAliases:       []string{},
			sqlTables:           []string{"chat"},
			sqlTablesAliases:    []string{},
			whereFields:         []string{"id"},
			whereAliases:        []string{},
			whereFieldsOperator: []string{EQUAL},
			whereLogicOperators: []string{},
		},
		&SQLModel{
			updSetFields:        []string{"alias", "create_at"},
			updSetAliases:       []string{"", "c"},
			sqlTables:           []string{"chat"},
			sqlTablesAliases:    []string{"c"},
			whereFields:         []string{"id", "author"},
			whereAliases:        []string{"c", ""},
			whereFieldsOperator: []string{EQUAL, EQUAL},
			whereLogicOperators: []string{AND},
		},
		&SQLModel{
			updSetFields:                    []string{"alias", "create_at", "client_hash"},
			updSetAliases:                   []string{"c", "c", "cu"},
			sqlTables:                       []string{"chat", "chat_user", "user"},
			sqlTablesAliases:                []string{"c", "cu", "u"},
			whereFields:                     []string{"id", "author"},
			whereAliases:                    []string{"c", "c"},
			whereFieldsOperator:             []string{EQUAL, EQUAL},
			whereLogicOperators:             []string{AND},
			whereDependenciesFields:         []string{"id", "id_chat", "id", "author"},
			whereDependenciesAliases:        []string{"c", "cu", "u", "c"},
			whereDependenciesOperator:       []string{EQUAL, EQUAL},
			whereDependenciesLogicOperators: []string{AND},
		},
	}

	for i := 0; i < len(params); i++ {

		model = params[i]

		model.PrepareQuery()

		if model.query != required[i] {
			t.Fatal("Failed generate Prepare Query to sql! \nRequired:", required[i], "; \nrezult:  ", model.query)
		}
	}
}

func TestPrepareInsertQuery(t *testing.T) {

	t.Log("Generate Prepare Insert Query to sql")

	var model *SQLModel

	required := []string{
		"INSERT INTO chat(author, alias, create_at, status) VALUES (?, ?, ?, ?)",
	}
	params := []*SQLModel{
		&SQLModel{
			insertFields:  []string{"author", "alias", "create_at", "status"},
			updSetAliases: []string{},
			sqlTables:     []string{"chat"},
		},
	}

	for i := 0; i < len(params); i++ {

		model = params[i]

		model.PrepareQuery()

		if model.query != required[i] {
			t.Fatal("Failed generate Prepare Query to sql! \nRequired:", required[i], "; \nrezult:  ", model.query)
		}
	}
}

func TestPrepareDeleteQuery(t *testing.T) {

	t.Log("Generate Prepare Delete Query to sql")

	var model *SQLModel

	required := []string{
		"DELETE FROM chat WHERE author=? AND (alias<>? OR create_at=?) AND id_chat<?",
	}
	params := []*SQLModel{
		&SQLModel{
			deleteSql:           true,
			sqlTables:           []string{"chat"},
			whereFields:         []string{"author", "alias", "create_at", "id_chat"},
			whereAliases:        []string{},
			whereFieldsOperator: []string{EQUAL, NOT_EQUAL, EQUAL, LESS},
			whereLogicOperators: []string{AND, OR, AND},
		},
	}

	for i := 0; i < len(params); i++ {

		model = params[i]

		model.PrepareQuery()

		if model.query != required[i] {
			t.Fatal("Failed generate Prepare Query to sql! \nRequired:", required[i], "; \nrezult:  ", model.query)
		}
	}
}
