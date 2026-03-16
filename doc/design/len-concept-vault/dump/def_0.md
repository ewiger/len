so we starts the discussion about a need of a high level spec language: "1. Natural language does not map uniquely into logic; the missing block is **semantic interpretation**.

2. For programming, the main failure is usually **frontend semantic hallucination**, not the logic backend.

3. Prolog is useful for executable rule fragments, but it only matches a restricted slice of logic well and often distorts semantics.

4. Lean is better as a **trusted verification kernel** because it supports explicit semantics, proof objects, and meta-reasoning.

5. Neither Prolog nor Lean solves ambiguity by itself; they only operate on whatever formalization they are given.

6. The right architecture is not “NL → Prolog/Lean,” but **input → semantic representation / IR → reasoning / verification backend**.

7. In programming, the frontend should not invent theorems directly; it should produce **structured candidate specs, assumptions, ambiguities, and proof obligations**.

8. A good system treats the LLM as an **untrusted proposal engine**, not as the final source of truth.

9. The verifier should decide correctness by checking that a candidate program **satisfies the formal specification**, not by plausibility or tests alone.

10. Under your idealized assumptions, the system becomes: **formal spec kernel → candidate synthesis/search → certified acceptance or rejection**.

11. The formal spec language for programs needs more than plain FOL; it should express **types, pre/postconditions, invariants, recursion, state, and refinement**.

12. The LLM’s best role is heuristic search: proposing **programs, invariants, lemmas, decompositions, and repairs**.

13. Fast rejection should come from **testing, symbolic execution, SAT/SMT, bounded model checking, and abstract interpretation**.

14. Final acceptance must require a **checked proof, refinement proof, or certificate** in a trusted kernel such as Lean.

15. The main operational loop is: **specify → propose candidate → refute quickly if possible → prove/refine → extract failure → revise → repeat**.

16. The right mathematical goal is: for specification `S` and candidate program `P`, prove **`M(P) ⊨ S`**.

17. It is often better to use **relational specifications** (`Spec(input, output)`) rather than only functional equalities.

18. Search should be constrained by a **grammar or IR**, not left as arbitrary raw code generation.

19. Even with perfect formal input, remaining risks include **wrong semantic embedding, incomplete specs, proof automation limits, and translation bugs**.

20. The deepest principle of the whole design is: **never trust generation; trust only checked preservation of semantics**.


