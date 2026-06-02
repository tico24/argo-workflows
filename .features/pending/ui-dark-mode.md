Description: Add an optional dark mode (system default, toggle, persisted)
Authors: [Tim Collins](https://github.com/tico24)
Component: General
Issues: 5037

Adds dark mode (closes #5037). Follows your OS setting by default, with a toggle to force light or dark that sticks via localStorage. The theme's applied to the root element before the app mounts so there's no white flash on load. Both themes pass AA contrast.
