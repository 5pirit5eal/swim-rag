---
description: Protocol for executing tasks from a task list, updating their status, documenting decisions, marking relevant files, and suggesting solution design refinements. Includes a visual workflow.
alwaysApply: false
---

# Protocol: Task Execution & Update

When a user indicates they are working on or have completed a task from a project task list, follow this protocol. This ensures tasks are properly updated, and relevant information is captured.

## Workflow Overview

The following Mermaid diagram illustrates the step-by-step process for task execution and updating the task list:

```mermaid
graph TD
    BA["Start Task Execution Cycle (TASK MODE)"] -- User indicates work on task --> BB[AI and User identify task list e.g. TASKS.md FEATURE_NAME.md]
    BB -- Task List --> BC[AI: Identify current state of the task list completion]
    BC -- Current Task --> BD{User: Confirm task selected from list}
    BD -- Wrong Task --> BC
    BD -- Correct Task --> BE[AI: Start implemeting specified task]
    BE --> BF["AI: Update Task File IMMEDIATELY and move task to 'In Progress' tasks"]
    BF --> BG[AI: Implement task]
    BG -- User approval --> BH["Mark Task as Completed Change to `-[x]` and move task from In Progress to Completed section"]
    BH --> BI[Discuss with User: Suggest Additions and Refinements to Solution Design]
    BI --> BI_Approval{"Add suggestions to implementation?"}
    BI_Approval -- Suggestions Approved --> BJ[AI: Update task list]
    BJ --> BK
    BI_Approval -- Suggestions Declined --> BK{New Tasks/Sub-Tasks identified during execution?}
    BK -- Yes --> BL["AI: Add new Tasks Sub-Tasks In Progress or Future Tasks"]
    BL --> BM[Confirm with User: Task List Updated Satisfactorily]
    BK -- No --> BM
    BM --> BN{Is task list complete?}
    BN -- Still open tasks --> BC
    BN -- Task list complete --> BS[End Task Update Cycle Await Next Action]
```

## General reminders

- **Standard Structure & Content:** A task list file has the following sections (`# [Feature Name] Implementation`,
   `## Completed Tasks`, `## In Progress Tasks`, `## Future Tasks`,
   `## Implementation Plan`, `### Relevant Files`). **When updating maintain these sections**.
- When updating the task list reference the CREATE TASK rule (`.clinerules/create-task.md`).
- Continue with the EPCC cycle after completing this EXECUTE TASK rule.
