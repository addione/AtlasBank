# Task Improvement Suggestions

## Original Task
"create separate files for models and optimize the migration process"

## How the Task Could Have Been More Specific

### 1. **Clearer Scope Definition**
Better task formulation:
```
"Refactor internal/database/models.go by:
1. Creating separate files for each model (User, OTP, Account, Transaction, Notification)
2. Moving models to a new internal/models package
3. Updating all import references across the codebase
4. Optimizing database migration files with better indexes"
```

**Why it's better:**
- Explicitly mentions the source file (`internal/database/models.go`)
- Lists specific models to separate
- Clarifies the target package structure
- Defines what "optimize" means (better indexes)

### 2. **Migration Optimization Criteria**
Better specification:
```
"Optimize migrations by:
- Adding composite indexes for frequently queried field combinations
- Adding indexes on status/type fields for filtering
- Consolidating redundant ALTER TABLE statements
- Ensuring migrations are idempotent (IF NOT EXISTS)"
```

**Why it's better:**
- Defines specific optimization techniques
- Provides measurable criteria
- Clarifies what "optimize" means in this context

### 3. **Expected Outcomes**
Better task formulation:
```
"Expected outcomes:
- Each model in its own file under internal/models/
- All services updated to use new import paths
- Build passes with no errors
- Migration files have optimized indexes for common queries
- Documentation of changes created"
```

**Why it's better:**
- Clear success criteria
- Testable outcomes
- Includes documentation requirement

### 4. **Constraints and Preferences**
Better specification:
```
"Constraints:
- Maintain backward compatibility where possible
- Don't break existing API contracts
- Keep migration files numbered sequentially
- Preserve all existing model fields and relationships"
```

**Why it's better:**
- Sets boundaries for the refactoring
- Prevents breaking changes
- Clarifies what should NOT be changed

### 5. **Performance Goals**
Better task formulation:
```
"Optimize migrations for:
- User verification queries (is_verified field)
- Account status filtering (active/frozen/closed)
- Transaction type and status lookups
- OTP validation queries (user_id + action + is_used)
- Notification status tracking per user"
```

**Why it's better:**
- Identifies specific query patterns to optimize
- Provides context for index decisions
- Helps prioritize optimization efforts

## Recommended Task Format Template

```markdown
### Task: [Clear, concise title]

**Context:**
- Current state: [What exists now]
- Problem: [What needs improvement]

**Objectives:**
1. [Specific goal 1]
2. [Specific goal 2]
3. [Specific goal 3]

**Scope:**
- In scope: [What should be changed]
- Out of scope: [What should NOT be changed]

**Success Criteria:**
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

**Constraints:**
- [Any limitations or requirements]

**Expected Deliverables:**
- [File/code changes]
- [Documentation]
- [Tests]
```

## Example: Improved Version of Your Task

```markdown
### Task: Refactor Database Models and Optimize Migrations

**Context:**
- Current state: All models are in a single file `internal/database/models.go`
- Problem: Hard to maintain, navigate, and scale as models grow

**Objectives:**
1. Separate each model (User, OTP, Account, Transaction, Notification) into individual files
2. Move models from `internal/database` to a new `internal/models` package
3. Update all import references in services, handlers, and database files
4. Optimize migration files with composite indexes for common query patterns

**Scope:**
- In scope:
  - Creating new model files in `internal/models/`
  - Updating imports in all service files
  - Adding performance indexes to migrations
  - Creating documentation of changes
  
- Out of scope:
  - Changing model field definitions
  - Modifying API contracts
  - Altering database schema structure

**Success Criteria:**
- [ ] Each model has its own file in `internal/models/`
- [ ] All imports updated from `database.*` to `models.*`
- [ ] Code builds successfully with no errors
- [ ] Migration files include optimized indexes for:
  - User verification queries
  - Account/transaction status filtering
  - OTP lookup by user+action+used
  - Notification queries by user+status
- [ ] Documentation created explaining changes

**Constraints:**
- Maintain backward compatibility
- Keep existing migration file numbering
- Don't modify existing model fields
- Ensure all tests still pass

**Expected Deliverables:**
- 5 new model files in `internal/models/`
- Updated service files with new imports
- Optimized migration files with new indexes
- MODELS_REFACTORING.md documentation
```

## Key Takeaways

1. **Be Specific**: Instead of "optimize", specify what optimizations (indexes, queries, etc.)
2. **Define Scope**: Clearly state what's in and out of scope
3. **Set Criteria**: Provide measurable success criteria
4. **Give Context**: Explain why the change is needed
5. **List Constraints**: Mention what should NOT be changed
6. **Specify Deliverables**: List expected outputs (code, docs, tests)

## What Worked Well in Your Task

✅ Clear action verbs ("create", "optimize")
✅ Identified the main components (models, migrations)
✅ Concise and to the point

## What Could Be Improved

❌ Didn't specify which models to separate
❌ Didn't define what "optimize" means
❌ Didn't mention the target package structure
❌ Didn't specify success criteria
❌ Didn't mention documentation needs
