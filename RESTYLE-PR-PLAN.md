# UI restyle — PR plan & status

Status doc for the UI restyle work, kept in git so the plan survives independent of any external notes. Last updated 2026-06-02.

## What's being shipped

Three stacked PRs, all built on top of PR #16159 (`ui-modernization`, the dependency-stack modernization, currently open against `argoproj/argo-workflows`):

1. **Restyle baseline** — visual-only modernization to match the Argo CD look (nav rail, consistent buttons, tidier filters panel, monospace log console, theme variables, AA contrast + focus rings). No behaviour changes.
2. **Filter sidebar → right + collapsible** — moves the filter panel from left to right and lets it collapse, so the table reclaims width and truncates fewer columns. CSS `order`, no DOM reorder. State persisted in localStorage. Applies to workflows / cron workflows / templates / reports.
3. **Optional dark mode** — follows the OS setting by default, with a toggle (light/dark) that persists; theme applied on the root element pre-mount to avoid a white flash. AA in both themes.

## Branches (on the tico24 fork)

```
ui-modernization              (#16159, base of everything)
 └ ui-restyle-argocd-faithful   PR A — restyle baseline
   └ ui-filter-move             PR B — filter sidebar right + collapsible
     └ ui-dark-mode             PR C — optional dark mode
```

Each commit: `feat(ui):` conventional title, DCO `Signed-off-by`, a `.features/pending/*.md` feature file. Strict isolation — PR B contains **zero** dark-mode code (dark mode may not be approved, so the filter move must stand alone).

## Where the PRs go

The real PRs target **`argoproj/argo-workflows`** (cross-repo from the tico24 fork), same as #16159.

**Blocked until #16159 merges.** `argoproj` has no `ui-modernization` branch (it only exists on the fork; #16159 is cross-repo onto `argoproj:main`), so a stacked PR can't be based on it upstream — GitHub rejects it ("Base ref must be a branch"). Once #16159 merges into `argoproj:main`:

- PR A → base `main` (clean restyle-only diff once #16159 is in main)
- PR B → base `ui-restyle-argocd-faithful`
- PR C → base `ui-filter-move`

(B and C bases must be pushed to the fork first; heads are `tico24:<branch>`.)

## Dress-rehearsal PRs (on the fork)

Drafts opened on `tico24/argo-workflows` to validate wiring and act as copy source for the real raise:

- A — fork PR #180 (base `ui-modernization`)
- B — fork PR #181 (base `ui-restyle-argocd-faithful`)
- C — fork PR #182 (base `ui-filter-move`)

Each uses the repo PR template (Motivation / Modifications + screenshots / Verification / Documentation / AI) with a generic AI-use disclosure (no tool names). Screenshots are hot-linked from the `pr-assets` branch on the fork so they render in the description without landing in any code diff.

## Issues to link

- **Dark mode** → existing `argoproj/argo-workflows#5037` ("UI: Dark mode"). Done.
- **Restyle** → no existing issue. Needs a new one filed on argoproj.
- **Filter move** → no exact match (#9578 is about collapsing groups *within* the table, not relocating the sidebar). Needs a new one filed on argoproj.

Draft issue text:

> **Modernize the UI styling to match Argo CD** — The Workflows UI looks dated next to Argo CD, despite both being built on the same component library. Proposes a visual-only pass: recoloured nav rail, consistent buttons, a tidier filters panel, monospace log console, and theme variables instead of hardcoded colours, plus fixing contrast and keyboard focus rings that fail AA in places. No behaviour changes.

> **Move the filter sidebar to the right and make it collapsible** — The filter panel sits on the left and eats horizontal space, so the table truncates columns. Proposes moving it to the right and making it collapsible (state persisted) so users reclaim table width. Distinct from #9578 (which is about collapsing groups within the table).

## Open items before the real raise

- [ ] File the two new argoproj issues (restyle, filter move); swap `Fixes #TODO` → real numbers in PR A/B bodies and the `Issues:` line of their feature files.
- [ ] Reconcile the sign-off email: commits are signed off as `tim@light.inc` but local git config `user.email` is `tim@thecollins.team`. Pick one so the DCO bot matches the GitHub account.
- [ ] Wait for #16159 to merge, then raise A/B/C on argoproj per the bases above (copy from the fork drafts).

## Notes

- Branch SHAs change every time a commit is amended — re-check before the real raise rather than trusting any recorded SHA.
- The shell here has a `GIT_EXTERNAL_DIFF` wrapper that compacts `git diff`; run verification greps via `env -u GIT_EXTERNAL_DIFF git ...`.
- Baseline restyle is also a recoverable rollback point on its own branch.
