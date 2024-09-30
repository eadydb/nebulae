package consts

var (
	ProjectName = "nebulae"
	ConfigName  = "nebulae.yaml"
)

var (
	Neo4jProject    = "Project"
	Neo4jPom        = "Pom"
	Neo4jMiddleware = "Middleware"
	Neo4jParent     = "Parent"
	Neo4jPlugin     = "Plugin"
	Neo4jModule     = "Module"
)

// CommonLanguages maps common programming languages to their file extensions or identifiers
var CommonLanguages = map[string][]string{
	"Python":      {".py"},
	"JavaScript":  {".js", ".jsx", ".mjs"},
	"TypeScript":  {".ts", ".tsx"},
	"Java":        {".java"},
	"C++":         {".cpp", ".cxx", ".cc", ".h", ".hpp"},
	"C":           {".c", ".h"},
	"Go":          {".go"},
	"Kotlin":      {".kt", ".kts"},
	"Rust":        {".rs"},
	"Scala":       {".scala"},
	"Lua":         {".lua"},
	"Shell":       {".sh", ".bash"},
	"SQL":         {".sql"},
	"R":           {".r", ".R"},
	"Groovy":      {".groovy"},
	"Objective-C": {".m", ".mm"},
	"Assembly":    {".asm", ".s"},
	"HTML":        {".html", ".htm"},
	"Markdown":    {".md", ".markdown"},
	"Maven":       {"pom.xml"},
	"Gradle":      {"build.gradle", "settings.gradle"},
}
