Description: Modernize the UI styling to match Argo CD
Authors: [Tim Collins](https://github.com/tico24)
Component: General
Issues: TODO

The Workflows UI looks dated next to Argo CD, despite both being built on the same component library. This is a visual-only pass to close that gap: recoloured nav rail, consistent buttons, a tidier filters panel, monospace log console, and theme variables instead of colours hardcoded everywhere. I also fixed contrast and keyboard focus rings that were failing AA in a few places. No behaviour changes — just CSS and a few small markup tweaks.
