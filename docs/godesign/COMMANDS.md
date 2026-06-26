# Commands
Commands are the primary way of interacting with the vector databases and document stores, and interacting with the workflows. It's also the entry point for running agentic workflows.

The primary commands are going to be interacting with the documents and the notes in some ways. Notes can be edited in the file system, but they might be updated by [[Workflows]] or other commands that convert or alter them in some way.

---

## `notes`

The `notes` command is the primary interface for managing notes.

### `notes import`

Imports text blobs or files into the notes [[vector store]].

**Flags:**

- `-i` / `--import` — trigger import
- `-f` / `--file` — path to the file being imported

When imported through a file, the content is semantically chunked and stored into the notes vector and note table.

### `notes query`

Takes in a query string. It attempts a general answer first, then uses that attempt to drive a [[vector embedding]] query.

**Flags:**

- `--limit` / `-l` — limit the number of results returned
- `--rank` / `-r` — enables reranking; takes the initial user query and the vector results and ranks them against an expected rank number, returning the top `n` vector matches from the `m` limit embedding query

**Output format:**

Results are returned in a table format with string abbreviations, consisting of:

- Document ID
- Title
- Matching text
- Tags

---

## `document`

The `document` command is more complex than `notes`. It has several subcommands.

### `document create`

A more performative action — creates a document, potentially from existing notes or other documents.

**Flags:**

- `--note-ids` — one or more note IDs (comma-separated) to reference when creating the document
- `--doc-ids` — one or more document IDs (comma-separated) to reference
- `--prompt` — a prompt string or, in the future, a prompt ID/name from a [[prompt repository]]
- `--template` — file path to a document template
- `--output` / `-o` — output file path (relative or absolute); defaults to stdout stream if omitted

**Behavior:**

Executes the given prompt (or a default prompt) with the reference text and the [[document vector search]] tool. Once processing is complete:

1. Creates the document file if an output path is provided
2. Automatically registers it in the document table
3. Runs embedding and updates the document vector table

> In the future, a [[prompt repository]] should allow `--prompt` to accept a prompt ID or unique name, enabling reusable prompt types (e.g. `document-prompt`, `note-prompt`).

### `document import`

Ingests an existing document file and updates both the document table and document vector table.

**Flags:**

- `--file` / `-f` — path to the file to import

### `document update`

Updates an existing registered document.

**Flags:**

- `--note-ids` — comma-separated note IDs to reference
- `--doc-ids` — comma-separated document IDs to reference
- `--prompt` — prompt string or prompt ID
- `--file` — file path of the document being updated
- `--doc-id` — ID of the document being updated (alternative to `--file`)
- `--create` — if the document is not already registered in the document table, this flag allows the command to create the entry rather than failing, and will also trigger ingestion and vectorization of the document

Without `--create`, the command does a pre-check and fails if the document is not already registered.

**Tag behavior:**

On update, existing tags are carried forward and then re-tagged based on the updated document content.

### `document query`

Functions practically identically to `notes query`, but operates on the documents [[vector store]].

**Flags:**

- `--limit` / `-l` — limit number of results
- `--rank` / `-r` — enable reranking against top `n` results

**Output format:** same table format as `notes query` — ID, title, matching text, tags.

---

### Document Summarization

When interacting with a document or multiple documents, the document(s) should be summarized. Optionally, a user query can be passed in to influence the summary. Summaries should also include tags for overarching topics.

> Creating, importing, or updating a document will update its tags. On update, existing tags are used as a baseline and then re-tagged based on the new content.

---

### Document Folder Structure

Document folders can be organized however the user prefers, but the general guidance is:

- **Internal documents** — for internal documentation purposes only
- **External documents** — externally facing; created by the agent system but potentially published or shared with teammates

The recommended sub-directory structure is:

```
documents/
  internal/
    category/
      subcategory/      # optional
        document.md
  external/
    category/
      subcategory/      # optional
        document.md
```

Categories and subcategories are optional and entirely up to the user.

---

## `report`

The `report` command generates structured reports based on notes and documents. Reports are output to a `reports/` folder by default, but the output path can be overridden. Reports do not need a document table record.

**Flags (shared):**

- `--type` — report type (`daily`, `weekly`, `quarterly`, `annual`)
- `--template` — optional path to a report template file
- `--output` — override the output path

---

### Report Format

Each report type should follow a consistent structure:

1. **Topics** — the most relevant or important tags, listed at the top
2. **Wins** — progress made on certain topics
3. **Areas to Improve** — suggestions for topics to explore, with reasoning on why

Templates can be provided per report type, but the above structure is the default format.

---

### `report --type daily`

Queries the notes and document tables for:

- Notes created within the current day
- Documents created or updated within the current day

Generates a daily report that:

- Provides a general count of documents and notes created
- Does a more detailed description based on the actual notes and documents worked on

### `report --type weekly`

Considers changes over the past week, with Monday as the first business day.

### `report --type quarterly`

Takes in the weekly reports from the quarter and generates a quarterly summary.

### `report --type annual`

Should use a larger model with a larger [[context window]]. Takes in all weekly reports and quarterly reports to produce:

- A quarter-by-quarter summary
- An overall annual summary

