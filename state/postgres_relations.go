package state

import (
	"strconv"

	"gopkg.in/guregu/null.v3"
)

type PostgresRelation struct {
	Oid                    Oid
	DatabaseOid            Oid
	SchemaName             string
	RelationName           string
	RelationType           string
	PersistenceType        string
	Columns                []PostgresColumn
	Indices                []PostgresIndex
	Constraints            []PostgresConstraint
	ViewDefinition         string
	Options                map[string]string
	HasOids                bool
	HasInheritanceChildren bool
	HasToast               bool
	FrozenXID              Xid
	MinimumMultixactXID    Xid
}

type PostgresColumn struct {
	RelationOid  Oid
	Name         string
	DataType     string
	DefaultValue null.String
	NotNull      bool
	Position     int32
}

type PostgresIndex struct {
	RelationOid   Oid
	IndexOid      Oid
	IndexType     string // Equivalent with pg_am.amname, e.g. "btree", "gist", "gin", "brin"
	Columns       []int32
	Name          string
	IsPrimary     bool
	IsUnique      bool
	IsValid       bool
	IndexDef      string
	ConstraintDef null.String
	Options       map[string]string
}

type PostgresConstraint struct {
	RelationOid       Oid     // The table this constraint is on
	Name              string  // Constraint name (not necessarily unique!)
	Type              string  // c = check constraint, f = foreign key constraint, p = primary key constraint, u = unique constraint, t = constraint trigger, x = exclusion constraint
	ConstraintDef     string  // Human-readable representation of the expression
	Columns           []int32 // If a table constraint (including foreign keys, but not constraint triggers), list of the constrained columns
	ForeignOid        Oid     // If a foreign key, the referenced table
	ForeignColumns    []int32 // If a foreign key, list of the referenced columns
	ForeignUpdateType string  // Foreign key update action code: a = no action, r = restrict, c = cascade, n = set null, d = set default
	ForeignDeleteType string  // Foreign key deletion action code: a = no action, r = restrict, c = cascade, n = set null, d = set default
	ForeignMatchType  string  // Foreign key match type: f = full, p = partial, s = simple
}

// Fillfactor - Returns the FILLFACTOR storage parameter set on the table, or the default (100)
func (r PostgresRelation) Fillfactor() int32 {
	fstr, exists := r.Options["fillfactor"]
	if exists {
		f, _ := strconv.Atoi(fstr)
		return int32(f)
	}
	return 100
}

// Fillfactor - Returns the FILLFACTOR storage parameter set on the index, or the default (100)
func (i PostgresIndex) Fillfactor() int32 {
	fstr, exists := i.Options["fillfactor"]
	if exists {
		f, _ := strconv.Atoi(fstr)
		return int32(f)
	}
	if i.IndexType == "btree" {
		return 90
	}
	return -1
}
