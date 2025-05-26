---
description: Main Software Engineering Process
alwaysApply: true
---
# Protocol: Integrated EPCC Workflow (with Task List Management)

When a user presents a significant problem, a new feature request, or a complex task, guide them and yourself (the AI) through the following integrated EPCC workflow. This workflow combines high-level strategic planning (EPCC) with detailed task management using established protocols.

## Integrated EPCC Workflow Visualized

```mermaid
flowchart TD
    A["Start: User Problem/Feature Request"] -- Switch to Plan Mode --> B("EXPLORE Phase")
    B -- User provides context --> B1["AI: Read Context (NO code)"]
    B1 -- Context Gathered --> C("PLAN Phase")
    C --> C1["User: Choose thinking depth to develop Plan"]
    C1 --> C2_Decision{"User: Does this plan require a detailed task breakdown or new feature task list?"}
    C2_Decision -- Yes --> C2_InvokeTaskCreate["AI: Invoke CREATE TASK rule (...) -> Create/Update Task List"]
    C2_Decision -- No --> C2_DirectPlan["AI: Formulate Plan in Chat"]
    C2_InvokeTaskCreate -- Task List Ready --> C3["User: Reviews Overall Plan & Detailed Task List"]
    C2_DirectPlan -- Plan formulated --> C3
    C3 -- Plan NOT OK --> B
    C3 -- Plan OK --> C3_Decision["AI: Document Overall Plan Approach"]
    C3_Decision -- Switch to Task Mode --> D("CODE Phase")
    D --> D1["AI: Implement Solution"]
    D1 -- If using task list --> D1_InvokeTaskExecute["Invoke 'EXECUTE TASK rule'"]
    D1_InvokeTaskExecute -- List completed --> E("COMMIT & DOCUMENT Phase")
    D1 -- Overall Code Implementation Complete --> E
    E -- "Switch to Plan Mode" --> E1["AI: Update MEMORY BANK documents"]
    E1 --> E2["AI: General Documentation Updates"]
    E2 -- User approves documentation --> E3["AI: Prompt User to commit code and request creation of PR"]
    E3 -- "Pull-Request Requested" --> E4["AI: Draft Pull Request"]
    E3 -- "Pull-Request NOT Requested" --> F
    E4 --> F["Task End: EPCC Cycle Completed"]
```

## AI Instructions for Each Phase

### A. Start: User Presents Problem/Feature Request

- Acknowledge the request.
- If complex, suggest: "This seems like a good fit for our integrated EPCC workflow, which helps us explore, plan in detail (including creating a task list if needed), code, and then commit with proper documentation. Shall we begin? Please switch to PLAN MODE"
- If agreed, initiate the **EXPLORE Phase**.

### B. EXPLORE Phase (Node B, B1 in diagram)

- **Objective:** Gather comprehensive context.
- **AI Action (B1):**
  - State: **"Let's start exploring. I won't write any code yet."**
  - Request relevant files, documentation, URLs, or codebase areas from the user.
  - Analyze provided context. Request clarification on areas with insufficient information.
  - Summarize understanding and clarify doubts.
  - Interact with the User until sufficient context was gathered to complete the request
- **Transition:** Once context is clear, proceed to the **PLAN Phase**.

### C. PLAN Phase (Node C, C1, C2, C2_Decision, C2_InvokeTaskCreate, C3 in diagram)

- **Objective:** Develop a strategic plan, and if necessary, a detailed task breakdown in a task list file.
- **AI Action (C1):**
  - Prompt the User: "Let's create a plan. How deeply should I 'think' about this (basic, thorough, deep)?"
    - basic: Plan with straightforward implementation and assumptions within best-practice
    - thorough: Complex plan with reflection on implementation steps, trade-offs and optimizations
    - deep: Consider 2-3 ways of implementation, compare the options regarding trade-offs, fit for the task and the amount of work. Challenge the plan and actively search for ways to optimize.
  - Formulate a high-level strategic plan.
- **Decision Point (C2_Decision):**
  - Assess with the user: "Does this plan warrant a detailed breakdown into a new or existing task list (e.g., for a new feature or multiple steps)?"
- **If YES (Invoke CREATE TASK rule - C2_InvokeTaskCreate):**
  - State: "Okay, let's detail this out. I'll now use our CREATE TASK rule (referencing `.clinerules/create-task.md`) to create/update the specific task list with all the necessary tasks, sub-tasks, dependencies, and initial implementation notes."
  - **Execute the `.clinerules/create-task.md` protocol fully.** This involves asking the user for the target task list file, defining tasks, handling dependencies, outlining the implementation plan *within that task list*, etc. The output is an updated/new `.md` task list file.
- **User Review (C3):**
  - Present the overall plan. If a task list was created/updated, present that as the detailed part of the plan.
  - **If Plan NOT OK:** Return to **EXPLORE Phase (B)** or refine plan/task list.
  - **If Plan OK:**
    - Document this overall strategic approach using the task list documents.
    - Prompt User to switch to TASK MODE.
    - Proceed to **CODE Phase**.

### D. CODE Phase (Node D, D1, D1_InvokeTaskExecute in diagram)

- **Objective:** Implement the solution based on the plan and the detailed task list (if one was created).
- **AI Action (D1):**
  - State: "Great, the plan is approved. I'll proceed with the implementation, following the overall strategy and working through the tasks in the task list if we created/updated one."
  - Begin coding.
- **Iterative Task Execution (D1_InvokeTaskExecute):**
  - **As each significant task or sub-task from the task list is being worked on or completed:**
    - State: "I'm now focusing on/just completed task: '[Task Description from list]'. I'll use our EXECUTE TASK rule (referencing `.clinerules/work-task.md`) to update its status and capture details."
    - **Execute the `.clinerules/work-task.md` for THAT specific task.** This includes marking it complete, documenting decisions/challenges *for that task* in the task list's Implementation Plan, noting relevant files *for that task*, and identifying any new sub-tasks emerging from it.
- **Transition:** When all planned coding and tasks from the list (for the current scope) are complete, inform the user and proceed to **COMMIT & DOCUMENT Phase**.

### E. COMMIT & DOCUMENT Phase (Node E, E1, E2, E3, E4 in diagram)

- **Objective:** Finalize code, update, and commit all code and relevant documentation including the task list and broader project docs.
- **AI Action (E1):**
  - Promt User to switch to PLAN Mode.
  - Begin final documentation updates.
- **MEMORY BANK Update (E1):**
  - State: "Let's integrate our completed task list into the MEMORY BANK documents using the EXECUTE TASK rule."
  - **Modify relevant documents mentioned in `.clinerules/memory.md`.** This might involve removing the task file and updating parts of the MEMORY BANK documents described in `.clinerules/memory.md`.
- **Update General Project Docs (E2):**
  - State: "Now that the main work is done and the MEMORY BANK is updated, lets update other general project documentation like API Documentation or Changelogs based on this feature/fix"
- **User Review (C3)**
  - Assist or perform updates as instructed.
- **AI Action (E3, E4):**
  - Confirm code with user
  - Prompt the user to commit the code using the git cli (e.g. `git add --all && git commit -m "feat: ..."`)
  - If applicable, prompt the user to create a Pull-Request
    - **If PR requested** Draft a PR summary
    - **If PR rejected** Crack a joke fitting the request.
- **Transition:** Transition to the final steps of the EPCC cycle

### F. End: EPCC Cycle Complete (Node F in diagram)

- Summarize the steps taken in this EPCC cycle.

This integrated protocol provides a robust framework, ensuring that strategic planning (EPCC) is seamlessly connected with granular task management. Remember to explicitly mention when you are switching to or referencing the sub-protocols (`.clinerules/create-task.md` and `.clinerules/work-task.md`).
