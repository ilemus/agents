# Workflow System Design

## Overview

Workflows should be highly customizable, but there should be general [[categories]]. The categories help people know which workflows to create, expand upon, or edit — and they also help logically track all the workflows.

There are gonna be various structures, but each one should be formatted the same way. There should be a general template for a workflow:

- A **description up front** — what is it for, what does it accomplish
- **Step N** with a title
  - Under each step, a subheading about what you have to do to accomplish it — 1, 2, 3 things, what data to retrieve, how to follow up
- Steps should be **easily identifiable** as either a **user action** or an **agentic action**

---

## Categories

> *(to be expanded — categories help route people to the right workflow and keep everything logically organized)*

---

## Workflow: Bug Resolution

**Description:** Given a [[stack trace]] and a description of the problem, follow up on the debugging and either provide an [[incident report]] or a [[bug fix]] (a bug fix is kind of included in an incident report — it's more about whether you're doing it formally or informally).

**Entry Points:**
- External: user-reported bug with stack trace / description
- Internal: found bug at XYZ, here's the log or the reproducible situation

---

### Step 1 — Capture Logs `[agentic]`

Find the full thread logs that are relevant. Highlight and capture the relevant logs.

---

### Step 2 — Find Related Code `[agentic]`

Find the related code in the codebase.

---

### Step 3 — Pull Reference Data `[agentic]`

Find related data input/output from the database, or reference data basically.

---

### Step 4 — Research Libraries & Docs `[agentic]`

Research the relevant libraries. Reference the documentation or schemas.

---

### Step 5 — Analyze the Bug `[user + agentic]`

Analysis is kind of a human action, but the model can augment it. The agent might provide: *"given all this information, here are three possible routes to consider"* — like maybe there's a thread collision, maybe there's malformed data, that kind of thing.

It might also be useful here to integrate with the [[vector database]] to search for previous similar problems — so part of the workflow might be: **look up previous problems around this code or with high similarity**.

---

### Step 6 — Draft Formal Incident Report `[agentic + user]` `[living document]`

> Draft this before writing any code. Update it throughout the workflow as new information comes in.

Uses a [[incident report template]]. Sections:

- **Overview** — brief high-level highlight
- **Affected** — is this a one-time user issue, does it affect other users, scope
- **Root Cause & Traceback** — how it got there, what happened because of XYZ, data query issues, how many modules are involved, external dependency issues
- **Resolution / Conclusion** — database update, schema rewrite, API schema change, adding type checking or stricter validation, etc. *(filled in progressively)*
- **Related Merge Requests** — reference section *(updated after implementation)*
- **Test Results** — *(added post-merge)*

---

### Step 7 — Draft Internal Code Document `[agentic + user]` `[living document]`

> Draft alongside or just after the formal report. Updated each iteration.

Document internally:
- What part of the code is related
- Why it's because of this type of data
- References for what needs to change and *why* — is there an external reference doc that can be included for more detail on what should actually be there?

Both the formal report and the internal code document together guide the actual code implementation. When the workflow loops, each iteration gets a new generation of the internal code doc (gen 1, gen 2, etc.) linked to the same thread.

---

### Step 8 — Update Jira Ticket `[user + agentic]`

At relevant points in the workflow, update the [[Jira]] ticket with the logs and relevant code (if necessary — whether it's an internal thing or external). Hints like *"you should add a comment update to the issue ticket"* should be surfaced.

- Optional if no external ticket system

---

### Step 9 — Implementation `[user + agentic]`

Use the formal report + internal code document to guide the code change.

---

### Step 10 — Review, Test, Merge `[user]`

- Is there test coverage? Did you have to update or rewrite a test?
- Generate the git commit message or merge message by referencing:
  - The formal document
  - The internal code document
  - The code diff (if possible)
- Get it reviewed and merged
- Update both the formal report and internal code doc with any changes that came out of review

---

### Step 11 — Post-Merge Verification `[agentic + user]`

- Run/record test results
- Update the formal document with test results
- Include an example of how the fix actually changed things

---

### Step 12 — Loop if Needed `[agentic + user]`

If results are not successful in production, loop back. Create a new generation of the internal code document. Both documents (first gen, second gen, etc.) should be linked to the same thread — so when you're iterating, you have:

- Formal document *(single, updated throughout)*
- Internal code document — gen 1, gen 2, ... *(each iteration linked)*

---

### Email Thread (Optional — for external/user-reported bugs)

There's an email sub-workflow alongside the main steps. Templates:

1. **Initial response** — *"Based on what you provided, I'm looking into XYZ. It might be because of XYZ."*
2. **Proposed fix update** — *"I have a proposed code change that I think might be the fix. It's expected to be released by [date]."*
3. **Post-merge check-in** — *"I made a code change and tested it — could you see if it works for you now?"*
4. **If fix didn't work** — *"That first code change didn't work. I'm going to look into something else and I'll update you."*
5. **Resolution / close** — *"Thank you for reaching out about the issue. I hope the fix is working for you. If you'd like more details, I have a document [XYZ]."*

> All emails should reference the [[Jira]] ticket.

There should also be a way to input **new information mid-workflow** — like if the issue got canceled, or there's new context, so that can be fed in and the workflow adjusts.

---

## Workflow: Design & Spec

**Description:** A longer-term workflow for software design. Starts from an initial spec or idea and expands it into a full spec document with supporting research, design decisions, and implementation planning.

**Entry Point:** An initial spec or idea — the main thing to capture here. There should be a solid paragraph capturing the thinking and decisions. If you have shorthand notes from while you were learning about a topic, you can import those and the agent helps expand them into something more complete.

Shorthand notes might include things like:
- *review more about XYZ topic*
- *reach out to this type of user*
- *look into this internal system*

These prompts help generate the initial description or spec.

---

### Step 1 — Expand Shorthand Notes `[agentic + user]`

Import shorthand/rough notes. The agent helps guide and expand them into a structured initial description.

---

### Step 2 — Draft Primary Spec Document `[agentic + user]` `[living document]`

The primary [[spec document]] — can get pretty large, but should be broken up. Uses the [[spec document template]].

> This is a living document. It's updated throughout the workflow as research comes in, decisions are made, and implementation progresses.

---

### Template: Spec Document
```
### [Feature / System Name]

#### High-Level Description

[What is this, what does it do, why does it exist]

#### High-Level Design

[Architecture overview, major components, how they relate]

---

### Sub-Modules / Sub-Pieces

#### [Module Name]

- **Description:** ...
- **Design:** ...
- **Design Choices:** [Key decisions made for this piece and why]

---

### Implementation Design

#### [Piece Name]

- **Approach:** ...
- **Why this approach:** ...
- **External references:** [docs, schemas, relevant prior art]

---

### Package & Tool Investigations

#### [Tool / Library Name]

- **Purpose:** What it's being evaluated for
- **Summary:** High-level capability overview
- **Feature Comparison:** [table or list — for these N features it does X, scores Y]
- **Detailed Notes:** [For these other features, here are the issues or caveats]
- **Decision:** Selected / Not selected / Pending — and why

#### Comparison Summary

[High-level comparison across all options evaluated]

---

### Open Questions

[Things that still need research, review, or decisions]

---

### Related Documents

- [[internal code document]]
- [[incident report]] _(if applicable)_
- External documentation links
```
---

### Step 3 — Package & Tool Investigation `[agentic + user]`

For each external tool or library being considered:
- Review what options exist
- Compare options — *for these five features it does X, scores Y*
- A high-level comparison table
- A more detailed report where relevant — *for these 10 features, there are issues with these things*

All of this goes into the [[spec document]] under the Package & Tool Investigations section. Design decisions and product reviews should be included — it's part of the reasoning, not just a side note.

---

### Step 4 — User / Stakeholder Outreach (Optional) `[user]`

Email templates for reaching out to users:

- **Soft approach** — inquiry/statement: *"How do you use this? What works for your situation?"* — open-ended, follow-up oriented
- **Hard approach** — decision-oriented: *"We're deciding to do this — can you see if that works for your situation?"*

Document who you're reaching out to, why, and the details of those conversations. This doesn't have to be a formal document — it's more like internal conversation notes and supporting references. Ideally, you're reviewing these before storing them (don't just dump raw chunks — interpret first, but quote the raw data too).

These notes can go into the [[vector store]] for later reference.

---

### Step 5 — Break Spec into Tasks `[agentic + user]`

Break the spec document down into tasks:
- Research XYZ package → find the necessary parts
- Implement those parts
- Test those parts / write test code
- Integrate those parts
- Full dev lifecycle for each piece of the spec

---

### Step 6 — Collaboration Tracking `[user]`

If other people are collaborating on parts of the implementation:
- Your role is less hands-on — more like receiving updates and documenting summaries of what they say
- If they produce documentation, copy it and store it as **external documentation**
- Keep a clear separation: *what I'm working on and what I know* vs. *what they're working on and what they know*
- External docs should also go into the [[vector store]] for searchability

---

## General Notes

### Vector Store Integration

Across multiple workflows, the [[vector store]] is useful for:
- Searching previous similar bugs or problems (high similarity)
- Storing and retrieving shorthand conversation notes
- Indexing external documentation from collaborators
- Linking past incidents to current ones

### Document Versioning & Threading

When a workflow loops (especially [[Workflow: Bug Resolution]]), documents should be linked together as a thread — first gen, second gen, etc. — so you always have the full history of attempts tied to the same issue.

### Personal / Social Notes (Separate System)

There should be a separate, **entirely private** system for the human side of interactions — not just about the work, but about the people you're interacting with. Things like:
- When interacting with this person, bring up this topic
- They prefer X or respond in Y way
- Their mood or communication style in a given context

This is a way to document and keep track of social context without internalizing it. It should be **fully segregated** from work notes — a separate process entirely, with its own interaction model.
