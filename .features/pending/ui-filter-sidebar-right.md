Description: Move the filter sidebar to the right and make it collapsible
Authors: [Tim Collins](https://github.com/tico24)
Component: General
Issues: TODO

The filter panel sits on the left and eats horizontal space, so the table gets cramped and truncates columns. This moves it to the right and lets you collapse it to get that width back when you don't need filters open. Done with CSS order rather than reordering the DOM, and the collapsed state is remembered. Applies to the workflows, cron workflows, templates, and reports views.
