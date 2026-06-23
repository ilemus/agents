# Agents & Benchmarks

A unified project designed to act as an agentic project workspace. This repository focuses on two primary goals: project planning/tracking (via a "reverse agent" approach) and LLM agentic benchmarking.

---

## 🎯 Project Goals

### 1. Agentic Project Planner & Tracker (The "Reverse Agent")
Typically, an AI agent is the worker executing tasks. In this system, the **worker is the human (acting as the agent)**, while the system acts as the coordinator. 
- **Systematic Prompting & Structuring**: Guides you through defining, designing, and executing your work based on your prior inputs.
- **Decomposition**: Automatically breaks down large, complex topics into structured, manageable subtasks.
- **Progress Tracking**: Documents and tracks tasks in a standardized, machine-readable/agent-friendly format.

### 2. LLM Agentic Benchmarking Suite
A lightweight, extensible framework to outline and benchmark LLMs on complex agentic tasks.
- **Aspects Tested**:
  - **Tool Calling**: Single-turn, multi-turn, and parallel tool invocation accuracy.
  - **Long Context Management**: Retrieval, reasoning over large context windows, and needle-in-a-haystack tasks.
  - **Summarization**: Quality, density, and fidelity of long-document synthesis.
  - **Thinking & Reasoning Modes**: Evaluation of chain-of-thought, self-correction, and planning capabilities.
- **Report Generation**: Rapid test execution and generation of consistent, structured benchmark reports.

---

## 📁 Repository Structure

```text
├── docs/                 # Design documents and architecture specifications
├── benchmarks/           # Code, configurations, and scripts for LLM evaluation
├── planner/              # Tooling and prompts for the project tracking/planning system
├── LICENSE               # Project license (Apache 2.0)
└── README.md             # Project overview (this file)
```

For detailed system designs, architectures, and design principles, please refer to the documents in the [docs](file:///home/lemusi/vscode/agents/docs) directory.

---

## 🚀 Getting Started

*(Detailed installation, usage instructions, and benchmarking commands will be added as implementation progresses.)*
