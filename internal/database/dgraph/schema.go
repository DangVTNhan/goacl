package dgraph

// Schema contains the Dgraph schema for the GoACL ReBAC system
const Schema = `
id: string @index(exact) .
email: string @index(exact) .
name: string @index(fulltext) .
created_at: datetime .
updated_at: datetime .
description: string @index(fulltext) .
namespace: string @index(exact) .
object_id: string @index(exact) .
relation: string @index(exact) .
user_id: string @index(exact) .
userset: string .
action: string @index(exact) .
resource_type: string @index(exact) .
resource_id: string @index(exact) .
timestamp: datetime @index(hour) .
token: string @index(exact) .
type: string @index(exact) .
effect: string @index(exact) .
details: string .
ip_address: string .
user_agent: string .
conditions: string .
permissions: [string] .
granted_at: datetime .
expires_at: datetime .
rewrite_rules: string .
member_of: [uid] .
owns: [uid] .
has_role: [uid] .
members: [uid] .
child_groups: [uid] .
parent_group: uid .
child_resources: [uid] .
owner: uid .
user: uid .
group: uid .
role: uid .
resource: uid .
granted_by: uid .
relations: [uid] .
actor: uid .
`

// InitialNamespaces contains the initial namespace configurations
var InitialNamespaces = []NamespaceConfigData{
	{
		Name: "documents",
		Relations: []RelationConfigData{
			{
				Name: "owner",
				RewriteRules: `{"union": {"child": [{"_this": {}}]}}`,
			},
			{
				Name: "editor",
				RewriteRules: `{"union": {"child": [{"_this": {}}, {"computed_userset": {"relation": "owner"}}]}}`,
			},
			{
				Name: "viewer",
				RewriteRules: `{"union": {"child": [{"_this": {}}, {"computed_userset": {"relation": "editor"}}]}}`,
			},
			{
				Name: "parent",
				RewriteRules: `{"union": {"child": [{"_this": {}}]}}`,
			},
		},
	},
	{
		Name: "folders",
		Relations: []RelationConfigData{
			{
				Name: "owner",
				RewriteRules: `{"union": {"child": [{"_this": {}}]}}`,
			},
			{
				Name: "editor",
				RewriteRules: `{"union": {"child": [{"_this": {}}, {"computed_userset": {"relation": "owner"}}]}}`,
			},
			{
				Name: "viewer",
				RewriteRules: `{"union": {"child": [{"_this": {}}, {"computed_userset": {"relation": "editor"}}, {"tuple_to_userset": {"tupleset": {"relation": "parent"}, "computed_userset": {"relation": "viewer"}}}]}}`,
			},
			{
				Name: "parent",
				RewriteRules: `{"union": {"child": [{"_this": {}}]}}`,
			},
		},
	},
	{
		Name: "groups",
		Relations: []RelationConfigData{
			{
				Name: "member",
				RewriteRules: `{"union": {"child": [{"_this": {}}, {"tuple_to_userset": {"tupleset": {"relation": "parent"}, "computed_userset": {"relation": "member"}}}]}}`,
			},
			{
				Name: "admin",
				RewriteRules: `{"union": {"child": [{"_this": {}}]}}`,
			},
			{
				Name: "parent",
				RewriteRules: `{"union": {"child": [{"_this": {}}]}}`,
			},
		},
	},
	{
		Name: "organizations",
		Relations: []RelationConfigData{
			{
				Name: "member",
				RewriteRules: `{"union": {"child": [{"_this": {}}]}}`,
			},
			{
				Name: "admin",
				RewriteRules: `{"union": {"child": [{"_this": {}}, {"computed_userset": {"relation": "owner"}}]}}`,
			},
			{
				Name: "owner",
				RewriteRules: `{"union": {"child": [{"_this": {}}]}}`,
			},
		},
	},
}

// NamespaceConfigData represents the structure for namespace configuration
type NamespaceConfigData struct {
	Name      string               `json:"name"`
	Relations []RelationConfigData `json:"relations"`
}

// RelationConfigData represents the structure for relation configuration
type RelationConfigData struct {
	Name         string `json:"name"`
	RewriteRules string `json:"rewrite_rules"`
}

// GetSchemaWithoutTypes returns the schema without type definitions for updates
func GetSchemaWithoutTypes() string {
	return `
# Indexes for performance optimization
id: string @index(exact) .
email: string @index(exact) .
name: string @index(fulltext) .
namespace: string @index(exact) .
object_id: string @index(exact) .
relation: string @index(exact) .
user_id: string @index(exact) .
action: string @index(exact) .
resource_type: string @index(exact) .
resource_id: string @index(exact) .
timestamp: datetime @index(hour) .
token: string @index(exact) .
type: string @index(exact) .
effect: string @index(exact) .

# Composite indexes for common query patterns
RelationTuple.namespace_object_relation: string @index(exact) .
RelationTuple.namespace_user_relation: string @index(exact) .
RelationTuple.object_relation: string @index(exact) .
AuditLog.resource_timestamp: string @index(exact) .
`
}
