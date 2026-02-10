# Stakeholder Review

Items annotated with task IDs. ✅ = resolved in Sprint 1/2, → T-NNN = tracked in Task Board.

- it seems saving does not work anymore. The autosave button shows an Error but we do not get feedback what does not work ✅ (Sprint 1 — save error tooltip + click-to-retry)
- an own button for saving or alternatively clicking on autosave to trigger the save would be helpful? ✅ (Sprint 1 — manual save button + saveNow callback)
- the sliding of the elements in the cv works but when you let go of the slider the element in the middle of the editor snaps to its new location. maybe it should rather also have a nice animation to make it more smooth ✅ (Sprint 2 — animateLayoutChanges + fallback transition)
- the cv itself looks really bland and it also lacks any color picker. We have added 3 templates in the data directory of this repository showing our main templates we want to use. Their accent color is the same in all, we might want to have our editor be able to use either of them and also set up the color ✅ (Sprint 2 — 3 template renderers + accent color picker)
- when the editor is only filling half the screen you can only make the middle part smaller but not the sidebars. This seems unresponsive ✅ (Sprint 2 — auto-collapse at breakpoints 768px/1024px)
- when I have a section I want to edit i should be able to also edit it by double clicking on it. Currently I can only edit with clicking on the button on the right sidebar ✅ (Sprint 2 — double-click selects section + opens properties)
- the additional information in education does not display, should it be visible in the cv? ✅ (Sprint 2 — all templates show location, description, grade)
- the pdf export does not work yet → T-003
- the word export does not work yet → T-004
- the json export actually works ✅ (already working)
- the editor lacks any way to returning to the dashboard ✅ (Sprint 1 — back-to-dashboard arrow in toolbar)
- the landing page itself is built very weirdly. It looks like a landing page for marketing which it should not. we should build it as an actual starting point like the dashboard. maybe we should merge the two ✅ (Sprint 1 — dashboard is now root page)
- the edit button on the dashboard gives a 404 ✅ (Sprint 1 — edit links to /editor)
- the burger menu on an existing cv does nothing ✅ (Sprint 1 — DropdownMenu with Edit/Delete)
- there is no way to change a template for a cv → T-001
