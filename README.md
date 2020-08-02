# len

recipe language for lazy programming

## Motivation

You can say it is a lazy language (hence where the name comes from). Reasonable extend of "being lazy" is a useful social trait, evident not only within animals but also among programs. The useful self-regulation properties of laziness can be compared to merits of altruism (such as open source) or implicit rationality (partially as turing completeness).

The goal of the programming language is brevity, with ability to reuse repeatable recipes and expand at need. As your program matures you may seek to address discovered imperfections, design limitations and optimization issues (things like CAP theorem). Programs evolve either by craft or by repetition (a.k.a *trial and error*). So you also simply need to move gradually with them, sometimes at your own speed of learnings. A lazy language can help with prototyping. The focus is on things happening at hand and at need, rather than ahead of the need. Only then you can enjoy or decently manage the complexity and the change rate of your own code.

Usually all complex things handled at need can become less complicated over time since somebody knows how to solve them, so there must be recipe for it.

## Design

Recipe-oriented language focusing on the following traits:

- transcriptional - transpile from shorter `.l` notations into `.py` targets linked with extensions
- conventional - certain naming conventions helps with linking extensions in C/Rust/Go
- meta - code-generation adds missing parts for you, so you can focus on contracts between modules, rather than concrete implementations and speed of the very first version
- versioned - elements of version control system are part of the language and code (module references, pins or git-trees are among first-citizens)


### Stories

Modules can be combined into multi-level categories. A collection of recipes can "tell" a story. For dependency management purposes, a complete collection of modules is a DAG. Meaning, that sublinking helps with re-using some exact parts from the upstreams, which can take place multiple times, but this is a basic git functionality.

The main story line increases with the complexity and sometimes becomes less stable. We consider the following categories:

- basic
- advanced
- experimental
- all

The last one is a a union of *all* of them, which is a complete category.


### Pins

Solution recipes are grouped into modules (primarily python once). Language does not focus on packages, rather arranges modules of your project repository into git trees.

Such collections or files allow pinning subtrees of remote repositories pointing at exact versions of the snippets. Snippets are a part of the recipe. This helps with tracking remote changes and merging the relevant update into your recipes. This also controls how the code-generation is using pins to pull the relevant changes in (at need).

