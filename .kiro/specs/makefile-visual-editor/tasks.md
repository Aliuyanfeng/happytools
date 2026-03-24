# Implementation Plan: Makefile Visual Editor

## Overview

Implement a visual Makefile editor module integrated into the HappyTools navigation. The feature includes a Go backend service for parsing/serializing Makefiles and managing persistent state (recent files, custom templates) via bbolt, and a Vue 3 + TypeScript frontend with a dependency graph view, target/variable editors, template library, and raw text mode.

## Tasks

- [x] 1. Navigation integration and module scaffolding
  - Add `makefileEditor` entry to `frontend/src/config/modules.ts` with id, nameKey, path `/makefile-editor`, icon `FileTextOutlined`, and theme `orange`
  - Add route `/makefile-editor` to `frontend/src/router/routes.ts` pointing to `MakefileEditor.vue`
  - Add i18n keys `home.modules.makefileEditor` and the `makefileEditor` namespace to `frontend/src/locales/zh-CN.ts` and `frontend/src/locales/en-US.ts`
  - Create skeleton `frontend/src/views/MakefileEditor/MakefileEditor.vue` with layout placeholder
  - _Requirements: 8.1, 8.2, 8.3_

- [x] 2. Go backend — data types and Makefile parser
  - Create `backend/services/makefile/types.go` defining `MakefileDoc`, `Target`, `Variable`, `RawBlock` structs with JSON tags
  - Create `backend/services/makefile/parser.go` implementing `Parse(content string) (*MakefileDoc, error)` that extracts Variables, Targets (name, deps, commands, isPhony), and preserves unrecognized blocks as RawBlock
  - Create `backend/services/makefile/printer.go` implementing `Print(doc *MakefileDoc) string` that serializes back to valid GNU Make syntax with Tab-indented commands and `.PHONY` declarations
  - _Requirements: 2.1, 2.2, 2.3, 2.5, 2.6_

  - [ ]* 2.1 Write property test for parser round-trip consistency
    - **Property 1: Round-trip property** — for any valid MakefileDoc, `Parse(Print(doc))` produces a semantically equivalent MakefileDoc
    - **Validates: Requirements 2.7**

  - [ ]* 2.2 Write unit tests for parser edge cases
    - Test `.PHONY` detection, multi-line commands, variable operators (`=`, `:=`, `?=`, `+=`), and unrecognized syntax preserved as RawBlock
    - _Requirements: 2.2, 2.3, 2.5_

- [x] 3. Go backend — bbolt store buckets and persistence helpers
  - Add `makefileRecentBucket` and `makefileTemplateBucket` bucket definitions to `backend/store/store.go` and register them in `Init`
  - Create `backend/store/makefile_options.go` with CRUD helpers: `SaveRecentFile`, `GetRecentFiles`, `SaveCustomTemplate`, `GetCustomTemplates`, `DeleteCustomTemplate`
  - _Requirements: 1.7, 6.3_

- [x] 4. Go backend — MakefileService implementation
  - Create `backend/services/makefile/makefileservice.go` implementing:
    - `OpenFile(path string) (*MakefileDoc, error)` — reads file, calls parser, records path in recent list (max 10)
    - `NewFile(dir string) (*MakefileDoc, error)` — creates empty doc, writes blank Makefile to disk
    - `NewFromTemplate(dir string, templateID string) (*MakefileDoc, error)` — initializes doc from built-in or custom template
    - `SaveFile(path string, doc *MakefileDoc) error` — serializes via printer, writes atomically (write to temp file then rename), returns I/O error without touching original on failure
    - `GetRecentFiles() ([]string, error)`
    - `GetTemplates() ([]Template, error)` — returns built-in + custom templates
    - `SaveCustomTemplate(name string, doc *MakefileDoc) error`
    - `DeleteCustomTemplate(id string) error`
    - `ValidateDependencies(doc *MakefileDoc) ([][]string, error)` — returns cycles as lists of target names
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 6.3, 6.4_

  - [ ]* 4.1 Write unit tests for SaveFile atomic write behavior
    - Test that I/O failure leaves original file unchanged
    - _Requirements: 1.6_

  - [ ]* 4.2 Write unit tests for ValidateDependencies cycle detection
    - Test direct cycle (A→B→A), transitive cycle (A→B→C→A), and acyclic graph
    - _Requirements: 3.5, 4.8_

- [x] 5. Go backend — built-in template definitions
  - Create `backend/services/makefile/templates.go` defining the five built-in templates as Go constants/structs:
    - Go multi-platform build (linux/windows/darwin)
    - Docker image build and push
    - Service deploy and restart
    - Frontend build (npm/yarn)
    - General clean
  - Each template includes name, description, preset Variables, and preset Targets with PhonyTarget declarations
  - _Requirements: 6.1, 6.2_

- [x] 6. Checkpoint — Ensure all backend tests pass
  - Ensure all Go tests pass, ask the user if questions arise.

- [x] 7. Wails bindings — register MakefileService
  - Register `MakefileService` in `main.go` (or wherever services are wired) so Wails generates TypeScript bindings
  - Verify generated bindings appear under `frontend/bindings/`
  - _Requirements: 1.1, 1.5, 4.2, 5.5_

- [x] 8. Frontend — Pinia store for editor state
  - Create `frontend/src/stores/makefileEditor.ts` with state: `currentDoc`, `currentPath`, `isDirty`, `recentFiles`, `templates`, `selectedTargetName`
  - Implement actions: `loadFile`, `newFile`, `newFromTemplate`, `saveFile`, `setSelectedTarget`, `addTarget`, `updateTarget`, `deleteTarget`, `addVariable`, `updateVariable`, `deleteVariable`, `mergeTemplate`
  - Cycle detection helper that calls backend `ValidateDependencies` and exposes `cycleWarnings`
  - _Requirements: 1.4, 3.3, 4.2, 4.5, 4.6, 4.8, 5.5_

- [x] 9. Frontend — file management panel component
  - Create `frontend/src/views/MakefileEditor/FilePanel.vue` with:
    - "Open File" button invoking Wails file dialog filtered to `Makefile`, `makefile`, `*.mk`
    - "New File" button with directory picker
    - "New from Template" button opening template selector modal
    - Recent files list (max 10) with click-to-open
    - Current file path display
    - "Save" button (disabled when not dirty)
  - Wire to Pinia store actions and display I/O error messages on failure
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7_

- [x] 10. Frontend — variable list and editor component
  - Create `frontend/src/views/MakefileEditor/VariableList.vue` displaying all Variables in a table (name, operator, value) with inline edit and delete
  - Create `frontend/src/views/MakefileEditor/VariableForm.vue` for add/edit with validation: name not empty (req 5.3), name not duplicate (req 5.4), operator selector (`=`, `:=`, `?=`, `+=`)
  - Wire save/delete to Pinia store actions which call backend SaveFile
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 11. Frontend — target list and editor component
  - Create `frontend/src/views/MakefileEditor/TargetForm.vue` with:
    - Target name input with validation: not empty (req 4.3), not duplicate (req 4.4)
    - PhonyTarget toggle switch
    - Dependency multi-select from existing target names; on selection change call `ValidateDependencies` and disable/warn on cycle (req 4.8)
    - Command multi-line textarea with auto Tab-indent on Enter key (req 4.7)
  - Create `frontend/src/views/MakefileEditor/TargetList.vue` listing all targets with edit/delete buttons
  - Wire to Pinia store; delete also removes from `.PHONY` list (req 4.6)
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 4.8_

- [x] 12. Frontend — dependency graph component
  - Create `frontend/src/views/MakefileEditor/DependencyGraph.vue` using a lightweight graph library (e.g., `@vue-flow/core` or `d3`) to render Target nodes and directed dependency edges
  - PhonyTargets rendered with distinct color/icon vs. regular Targets (req 3.2)
  - Nodes with cycle involvement highlighted in red with warning tooltip (req 3.5)
  - Support pan and zoom (req 3.4)
  - On node click emit selected target name to parent; parent updates Pinia `selectedTargetName` to show detail panel (req 3.3)
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

- [x] 13. Frontend — raw text editor mode
  - Create `frontend/src/views/MakefileEditor/RawEditor.vue` using a code editor component (e.g., `@codemirror/lang-makefile` or Monaco) with Makefile syntax highlighting
  - Bind editor content to `currentDoc` serialized text; on switch back to visual mode re-parse and refresh store (req 7.2)
  - Unrecognized syntax preserved as RawBlock in store, shown as read-only "raw content" block in visual mode (req 7.3)
  - Visual mode changes update raw editor content reactively (req 7.4)
  - _Requirements: 7.1, 7.2, 7.3, 7.4_

- [x] 14. Frontend — template library modal
  - Create `frontend/src/views/MakefileEditor/TemplateModal.vue` listing built-in and custom templates
  - Show template preview (serialized Makefile text) before applying (req 6.4)
  - "Apply (replace)" and "Merge" actions; merge only appends non-existing Targets and Variables (req 6.5)
  - "Save current as template" button with name input, persists via backend `SaveCustomTemplate` (req 6.3)
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 15. Frontend — assemble MakefileEditor main view
  - Assemble `MakefileEditor.vue` with:
    - Left sidebar: `FilePanel` + `VariableList`
    - Center: `DependencyGraph` (visual mode) or `RawEditor` (raw mode), with mode toggle button
    - Right panel: `TargetForm` (shown when a target is selected or being added), `TargetList`
    - `TemplateModal` triggered from `FilePanel`
  - Ensure visual ↔ raw mode toggle re-parses and syncs state correctly
  - _Requirements: 3.1, 4.1, 5.1, 7.1, 7.4_

- [x] 16. Final checkpoint — Ensure all tests pass
  - Ensure all Go and frontend tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- The Go parser must preserve unrecognized syntax as RawBlock to satisfy requirement 2.5 / 7.3
- Atomic file writes (write-to-temp + rename) are required by requirement 1.6
- Cycle detection must run client-side on every dependency selection change (req 4.8) and be visualized in the graph (req 3.5)
