# Project Customization

Run this after initial setup to configure the project context.

---

Scan this project and update `CONTEXT.md`:

1. **Detect the tech stack**
   - Language and version
   - Framework(s)
   - Database (if any)
   - Testing framework
   - CI/CD setup

2. **Map the project structure**
   - Key directories and their purpose
   - Entry points
   - Configuration files

3. **Identify conventions**
   - Code style (formatter, linter)
   - Existing patterns in the codebase
   - Git workflow (if .git exists)

4. **Document common commands**
   - How to run in development
   - How to run tests
   - How to build
   - Any project-specific scripts

5. **Note environment setup**
   - Required environment variables
   - Local setup steps

6. **Add first WORKLOG.md entry**
   - Session focus: "Initial setup and project familiarization"
   - Notes on what you discovered
   - Suggested next steps

**Important**: Do NOT modify `.claude/agents/` - they are designed to work with any project. All project-specific context goes in `CONTEXT.md`.