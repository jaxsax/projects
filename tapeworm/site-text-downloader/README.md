# Experiments with downloading article text from sites

We can't rely on pure html parsing, too tedious and too many variations.

We can use a mixture of techniques

- https://ahrefs.com/blog/open-graph-meta-tags/
- Firefox Reader View

Enhancers can combine data from different sources falling back to the shitty version only when absolutely necessary

What is a good common data structure we can build this upon? Or should it be a generally schemaless thing and when we need to implement search, we can fire off N requests to different indexes and merge the data together. Would that provide a generally useful search?
