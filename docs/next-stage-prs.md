# Next-stage PR plan — BiliOpinion hardening

Living checklist. One PR at a time, TDD, 8-way review, squash-merge, then check the box.

Parent goal acceptance (must all hold on `main` after PR6):

1. Progress URL never hangs blank; completed → report, failed → error, processing → SSE. Parse/confirm/settings never navigate to `/progress/undefined`. Video analyze cannot double-start. Restart does not duplicate history or re-run scrape/AI.
2. Score/radar use 1–10. Overview shows recommended brand + reason first. Dimension filter changes charts. Budget/scenario/special needs appear in recommendation when parse produced them. Confirm can add/remove brands, dimensions, keywords. Video progress labels match the video pipeline.
3. Analyzed comments persist. Brand score opens contributing original comments. History can re-analyze from store without Bilibili scrape.
4. `GET /api/config` never returns full AI key or full cookie. Server listens on `127.0.0.1`. Frontend has no hardcoded `http://localhost:8080`.
5. Playwright (mocked) covers 需求 → 确认 → 进度 completed → 报告. `/settings` page-only suite is not the sole e2e. Each PR has tests that drive shipped functions.

Non-goals: SaaS/login/share, full single-video product, default-on brand discovery UI, process-wide Bilibili token bucket, PDF/Excel rewrite, non-loopback bind, real auth.

---

## Child-goal template (copy into every PR)

```
TDD seams: write a failing test at the public seam first, then minimum code.
After green: 8 parallel reviews — adversarial, spec-consistency, correctness,
architecture, performance, testing, maintainability, compatibility.
Fix must-fix (bugs + spec breaks). Then squash-merge. Check the PR box below.
Do not start PR N+1 until PR N is merged and this file is updated.
```

8 review types:

| Type | Ask |
|---|---|
| Adversarial | How can a user or attacker break this? |
| Spec-consistency | Does the PR meet its slice of the parent acceptance criteria? |
| Correctness | Edge cases, races, wrong status, silent failure? |
| Architecture | Right module, no new coupling, no duplicate pipeline? |
| Performance | Extra AI/Bilibili/DB cost? |
| Testing | Tests drive shipped functions? Red-green? No tautology? |
| Maintainability | Names, size, dead code? |
| Compatibility | Existing happy path still works? Settings/history/report? |

---

## PR1 — `fix/local-security-and-recovery`

**Outcome:** Listen `127.0.0.1`. Mask config secrets. Recovery marks `processing` failed. No second history row. No full re-run.

**TDD seams:** `MaskSecret`; `GET /api/config` handler; listen address constant; `RecoverIncompleteTasks` against a real SQLite history table.

**Acceptance:**
- GET config returns empty or masked key/cookie, never the full stored secret.
- Empty secret on save does not overwrite an existing secret.
- `r.Run` / listen address is `127.0.0.1:8080`.
- Startup recovery does not call Execute and does not `Create` another `AnalysisHistory` for the same `task_id`.
- Timed-out (and leftover) `processing` rows become `failed`.

**Files:** `backend/main.go`, `backend/api/config.go`, `backend/task/recovery.go`, new tests.

- [x] PR1 merged

---

## PR2 — `fix/progress-and-confirm-hemostasis`

**Outcome:** Progress restores from history. Parse/confirm/settings honor errors. Video analyze submitting lock. Video step labels.

**TDD seams:** progress restore mapping helper; confirm/parse error → no progress navigation; video submitting flag.

**Acceptance:**
- Opening `/progress/:id` for completed task → report; failed → error message; processing → SSE.
- Non-OK parse/confirm or missing `task_id` never navigates to `/progress/undefined`.
- Settings save failure is shown as failure.
- Video analyze cannot start two tasks from a double click.
- Video progress steps are not “搜索相关视频”.

**Files:** `frontend/src/pages/Progress.tsx`, `Confirm.tsx`, `VideoConfirm.tsx`, `Settings.tsx` / `SettingsModal.tsx`, `hooks/useSSE.ts`, backend history lookup if needed.

- [ ] PR2 merged

---

## PR3 — `fix/report-scale-and-decision`

**Outcome:** 1–10 colors. Radar max 10. Overview recommendation first. Dimension filter wired. Hidden tabs do not mount all ECharts.

**TDD seams:** score-tone function (8.5 is success, not 90-based red); radar max constant 10; filter output actually passed into chart props.

**Acceptance:**
- Brand card / detail modal use 1–10 thresholds (≥8 green, ≥6 blue, ≥4 amber, else red).
- Radar axis max is 10.
- Overview shows recommended brand + short reason above the card wall.
- Changing dimension filter changes chart inputs.
- Inactive report tabs do not mount ECharts.

**Files:** `frontend/src/components/Report/**`, `frontend/src/pages/Report.tsx`.

- [ ] PR3 merged

---

## PR4 — `feat/editable-confirm-and-api-client`

**Outcome:** Editable brands/dimensions/keywords. Recommendation prompt includes budget/scenario/special needs. Single API client + Vite `/api` proxy. Dead `/settings` page deleted or routed.

**TDD seams:** confirm payload includes user-edited lists; recommendation generator / prompt includes budget/scenario/special needs when present.

**Acceptance:**
- Confirm page can add/remove brands, dimensions, keywords; those values are what `/api/confirm` receives.
- When parse produced budget/scenario/special needs, purchase recommendation text includes them.
- Frontend source has no `http://localhost:8080`. Calls go through one client / relative `/api`.
- Vite proxies `/api` to backend.

**Files:** `frontend/src/pages/Confirm.tsx`, `frontend/src/api/*`, `frontend/vite.config.ts`, `backend/report/generator.go`, `backend/ai/*`, `frontend/src/pages/Settings.tsx`.

- [ ] PR4 merged

---

## PR5 — `feat/comment-persist-and-reanalyze`

**Outcome:** Persist comments used for analysis. Drill-in to original text. Re-analyze from store without scrape.

**TDD seams:** persist + reload comments by history; re-analyze path does not call Bilibili scrape entry.

**Acceptance:**
- Comments that entered analysis are written to `raw_comments` (or equivalent).
- Brand score / detail can show contributing original comment text.
- History “re-analyze” rebuilds the report from stored comments without calling scrape.

**Files:** `backend/models/raw_comments.go`, `backend/task/executor.go`, `backend/api/history.go`, `backend/api` reanalyze handler, report UI drill-in.

- [ ] PR5 merged

---

## PR6 — `test/e2e-main-path`

**Outcome:** Mocked e2e for 需求 → 确认 → 进度 → 报告. `/settings`-only suite is not the sole e2e. Verification captures written.

**Acceptance:**
- Playwright (or equivalent) reaches a report screen with mocked parse/confirm/SSE/report.
- Old `/settings` page-only suite is not the only e2e.
- Scratch captures: go-test, config launch ×2, fe-unit, e2e.

**Files:** `frontend/e2e/*`, frontend unit test runner if added in earlier PRs.

- [ ] PR6 merged

---

## Final

- [ ] All five parent acceptance criteria hold on `main`
- [ ] This file’s PR1–PR6 boxes are checked
- [ ] Review notes exist under the goal scratch dir as `pr-N-reviews.md`
