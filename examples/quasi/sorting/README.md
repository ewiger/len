# sorting quasi examples

This folder is a small corpus for the accepted `ProceduralAlgorithm` quasi
style.

- `sorting.l1` holds the shared domain vocabulary and generic sorting laws
- `bubble_sort.l1`, `insertion_sort.l1`, `merge_sorting.l1`, and
  `quick_sort.l1` hold algorithm-specific relations, specs, and functions

The files are intended to validate together as one example corpus.

The current accepted direction is:

- keep algorithm linkage on `fn` with `implements`
- keep semantic laws in `spec`
- keep the host parser shallow and validate quasi lines by style profile