package db_schema

type RelationshipType string

const (
	HasOne     RelationshipType = "has_one"
	HasMany    RelationshipType = "has_many"
	BelongsTo  RelationshipType = "belongs_to"
	ManyToMany RelationshipType = "many_to_many"
	has        RelationshipType = "has"
)

type Relationships struct {
	HasOne            []*Relationship
	BelongsTo         []*Relationship
	HasMany           []*Relationship
	ManyToMany        []*Relationship
	Relations         map[string]*Relationship
	EmbeddedRelations map[string]*Relationships
	//TODO: нужен ли тут Mux sync.RWMutex
}

type Relationship struct {
	Name        string
	Type        RelationshipType
	Field       *Field
	References  []*Reference
	Schema      *Schema
	FieldSchema *Schema
	JoinTable   *Schema
}

type Reference struct {
	PrimaryKey    *Field
	PrimaryValue  string
	ForeignKey    *Field
	OwnPrimaryKey bool
}
