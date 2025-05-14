# Decision Log: Go Product Search API

---

### **1\. Problem Understanding**

The task is to build a backend API in Go that:

- Simulates 1 million in-memory product records
- Indexes them using full-text search (Bleve)
- Exposes a REST API endpoint `/search?q=...`
- Supports a graceful shutdown mechanism
- Uses Chi for HTTP routing

The goal is not just to deliver working code, but to showcase the decision-making and problem-solving process in an unfamiliar tech stack.

---

### **2\. Design & Implementation Decisions**

| **Area** | **Decision** | **Reasoning** |
| --- | --- | --- |
| Language | Go  | Required by task. Provides strong performance for concurrent applications. |
| HTTP Router | Chi | Lightweight, idiomatic Go HTTP router. Simple to integrate and well-maintained. |
| Search Library | Bleve (MemOnly) | Avoids disk I/O, fast indexing/searching in-memory. Ideal for quick prototyping. |
| Product Data | Generated 1M fake products with `fmt.Sprintf("Product %d")` | Avoided external dependencies; simple to scale for 1M records. |
| Categories | Static slice of categories | Adds slight realism for diversity in data. |
| Index Mapping | Default Bleve index mapping | Sufficient for simple string search needs. |
| Search Query | `bleve.NewQueryStringQuery(query)` | Simple full-text query for product name matching. |
| Results Cap | Limited to 50 matches | Keeps API responses performant and predictable. |
| Graceful Shutdown | Signal catching + `context.WithTimeout` | Clean termination and shutdown of the HTTP server. |
| Project Structure | Single main.go file | Keeps things minimalistic and focused, ideal for a timed take-home task. |
| go.mod | Initialized using `go mod init` | Required to manage dependencies using Go Modules. Avoided GOPATH complexity. |

---

### **3\. Technical Steps Taken**

- Used `rand.Seed` and `rand.Intn()` to generate randomized categories.
- Indexed each product using its ID as the key.
- Routed `/search` via Chi, parsed `q` param.
- Queried Bleve index and mapped search hits back to product slice.
- Implemented `os/signal.Notify` with `syscall.SIGINT/SIGTERM` to trigger shutdown.
- Initialized go.mod properly using `go mod init github.com/sppidy/go-product-search-api`.

---

### **4\. External Resources Consulted**

- [Chi router documentation](https://go-chi.io/#/pages/getting_started)
- [Bleve search documentation](https://blevesearch.com/docs/Home/)
- Go Blog and Modules Docs (for understanding go.mod)

---

### **5\. AI Prompts Used**

- "How to create in-memory Bleve index in Go?"
- "Chi router GET route example"
- "Go graceful shutdown example using context"
- "Generate fake products in Go"
- "How to index and search with Bleve"
- "How to initialize go.mod file in Go?"
- "Best practices for organizing Go web APIs"
- "Why use Bleve over ElasticSearch for local projects"
- "Go module vs GOPATH explanation"
- "Handle SIGINT and SIGTERM in Go"
- "Why go.mod is needed and how to fix go: error reading go.mod"
- "Should I use disk-based or memory-only Bleve index for fast search?"
- "Memory usage of Bleve with 1 million in-memory documents?"
- "How to optimize Bleve indexing for large data sets in Go"
- “Different types of indexing in Bleve“

---

### **6\. Learnings & Reflections**

- Bleve is a solid tool for lightweight full-text search in Go.
- Chi router is intuitive and clean for REST APIs.
- Graceful shutdown is crucial for real-world production services.
- Generating large datasets in-memory can be surprisingly fast and simple in Go.
- Learning to use Go modules was essential for managing third-party packages like Chi and Bleve.

---

### **7\. System Resource Usage (Observed)**

When generating and indexing 1 million products entirely in-memory:

| **Metric** | **Observation** |
| --- | --- |
| RAM Usage | ~16 GB peak usage (due to full in-memory Bleve index) |
| CPU Usage (search) | ~20–40% during high query load |
| API Response Time | ~10–60ms for top 50 results |
| Index Build Time | ~20–25 seconds to index all 1M records |
| Concurrency Handling | Chi handled concurrent requests without blocking |
| CPU Usage (While Generating) | ~80-100% |

**Tools Used**: `btop`, `time`, `curl`, and `hey` for basic load simulation.

**Tradeoffs & Notes**:

- Memory usage was significantly high due to Bleve’s in-memory structure.
- Would consider disk-based or hybrid approach for production scenarios.
- Pagination, query caching, or category-based filtering could reduce resource load.

---

**Prepared By:** Ramshouriesh R