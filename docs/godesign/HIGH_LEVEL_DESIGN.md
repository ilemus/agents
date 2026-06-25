## Overview

The project is designed around **[[Workflows]]** that accomplish different types of planning tasks.

---

## Workflows

### Task Splitting Workflow

The user enters a block of text or a reference document, and the agent identifies the main tasks and breaks them into smaller pieces.

The pieces are systematic — similar to a software design flow, but adaptable to any type of flow. The user has control over how the workflow is designed.

The user can specify the general workflow in a text document (likely Markdown). It highlights the general steps, but does **not** require strict adherence — each use case may follow a different branch or step, which is left to the LLM to determine.

The workflow agent references a **[[Step Definition File]]** provided by the user. There will likely be one or two example workflows included by default — probably a general system, coding, or project workflow.

---

### Workflow Steps (Example)

Steps might include:

1. **Research Phase** — Identifying topics to research or references to include
2. **Spec Expansion** — Contacting people for details or expanding on the [[Spec Document]]
3. **Spec Definition** — Defining the [[Spec Document]]

The [[Spec Document]] step includes sub-tasks, and a **general format/template** should be provided to the LLM for reference. This is not a full implementation — just a template.

**Example [[Templates]]:**
- Resolving an issue
- Researching a topic or comparing services/products
- Implementing a code piece or code-related task

---

### Post-Spec Workflow

After the [[Spec Document]] is created, it gets ingested to drive the next phase of the workflow, which may include:

- Looking up relevant dependencies (for software projects)
- Further research or next steps based on the [[Spec Document]]

---

### Optional Steps

Optional workflow steps may include things like:

- Emailing or contacting other people or teams
- Collecting more information

**[[Email Templates]]** should follow a general layout, e.g.:
- Reading task
- Reference or background
- A structured system for writing emails effectively

---

### User Input & Templating

Most output should ultimately come from the **user**. The agent should not auto-fill things, but should help structure or suggest content and allow the user to expand.

Template placeholders might use double-brace syntax, e.g. `{{task objective}}`, which references a defined value for the current task.

---

## Agent

### Workspace

The agent operates within a **[[Workspace]]** and has access to reference documents. Both should be configurable.

A default install should include:
- A general folder structure
- A dedicated [[Workspace]] folder per project

The [[Workspace]] stores structured data, including:
- All planned steps
- Documents written into the workspace folder

The agent should have **limited scope** — it should only be able to write within an allowed folder. It should have tools to read the existing [[Workspace]], know which files it has created, and be able to reference those when needed.

---

### Implementation

The implementation is planned as a **custom harness** — defining steps and planning logic — though an existing harness may be used in production.

Reference documents (especially the [[Step Definition File]]) should be ingested so that certain behaviors are more deterministic.

The agent should be able to **document timestamps** of various events.

---

### Interface

The agent is a **command line tool**. The user interacts with it via a single command, and can:

- Attach a document
- Reference a document already in the [[Workspace]]
- Notify it of updates (e.g. "I've updated the document in this workspace")

---

### Step Management & Next Steps

The agent should be able to suggest next steps, including:

- Multiple steps at once
- A hierarchy of steps (e.g. research phase, emailing phase)
- Suggestions based on where the user has made progress across multiple parallel tasks

The user may override, reject, or skip steps. There should be a way to manage the workflow accordingly.

[[Workflows]] should support **alternative routes**, with a good way to define fast vs. slow paths.

---

### Project & Task Management

The agent should function as a **[[Task Manager]]**, where each task has an entry in a database with its current status — making it queryable for next steps.

---

### Proofread Mode

A **proofread mode** should allow the agent to read documents the user has created and suggest updates.

---

### Code Project Support

For coding projects, the agent should optionally be able to:

- View the project and related code changes (diffs)
- Interpret those diffs in the context of the user's reference documents and created docs


# Workflow Knowledge & Storage Design

## Overview

Knowledge is a core part of how the system operates — the outputs of each step, whether document findings, communications, or other artifacts, should be stored. There are a couple of types to account for: structured files (markdowns) and shorthand notes. All of these feed into a [[Vector Database]] which stores semantic chunks of the notes and markdowns.

---

## Storage Categories

There should be three distinct categories, each stored in its own collection/table:

- **`workflow`** — task steps, state, and related operational data
- **`documentation`** — formatted markdown documents
- **`unstructured_notes`** — thought dumps and shorthand notes

Keeping these segregated makes it easier to query the right type of content at the right time, and avoids mixing structured task state with looser exploratory notes.

---

## Vector Store

The [[Vector Database]] stores semantic chunks with references back to the parent document. This parent-child chunking provides context throughout a project and makes things easily able to be referenced in the future.

**Example use case:** a workflow like "create a quarterly report" should be able to query across stored documents and notes — pulling highlights, topics discussed, findings — based on what was produced during that period.

### Metadata

Each record in the vector store should include:

- **Timestamp of ingestion** — when the record was stored
- **Creation date** — for markdown documents, when the document was first created
- **Modification history** — a separate table that records every time a document was updated (datetime stamps per update)

> Full versioning is deferred for now — don't want to store the same content 10 times over. Efficient versioning could be a future add-on.

---

## Retrieval

Agents should be able to query — or load in some structured way — the relevant docs based on two signals:

1. **The initial trigger/prompt** — what kicked off the workflow (email, chat message, log entry, stack trace + description, etc.)
2. **The current task state** — the refined, more detailed prompt as the task develops

So retrieval is context-sensitive: it starts broad based on the first prompt and narrows as the task gets more defined.

---

## Workflow Triggers

A workflow might be kicked off by any of the following:

- An email
- A chat message
- A log entry (e.g. a stack trace or error with a brief description of the request)

These triggers eventually get expanded upon as the workflow progresses.

---

## Workflow State Management

The [[Vector Database]] is separate from [[Workflow State Management]], which should be its own, more rigid database structure. That structure should include:

- `step` — the current step
- `step_name`
- `step_details`
- A graph-like tree state — with fields like `next_step` to traverse the workflow

This keeps operational state management cleanly separated from the knowledge/retrieval layer.