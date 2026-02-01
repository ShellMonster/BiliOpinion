# Issues & Problems - fix-home-input

## Known Issues
(To be populated as we discover problems)

## Pending Integration Issue
- **Confirm Page Mismatch**: `Home.tsx` now navigates to `/confirm?requirement=...`, but `Confirm.tsx` is still designed to read `url` parameter and display video info. The `Confirm` page needs to be refactored to handle the text requirement input and display appropriate analysis context (e.g., parsing the requirement instead of fetching video metadata).
